package envoy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	als "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	alf "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/envoyproxy/go-control-plane/pkg/test"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	google_protobuf1 "github.com/gogo/protobuf/types"
	"github.com/rancher/rio/pkg/settings"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/rancher/types/apis/core/v1"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

var (
	all                = "_all_"
	rioServiceSelector = labels.NewSelector()
)

func init() {
	l, err := labels.NewRequirement("rio.cattle.io", selection.Equals, []string{"true"})
	if err != nil {
		panic(err)
	}

	rioServiceSelector = rioServiceSelector.Add(*l)
}

func Register(ctx context.Context, context *types.Context) {
	c := &controller{
		nodeLister:        context.Core.Nodes("").Controller().Lister(),
		serviceController: context.Rio.Services("").Controller(),
		k8sServiceLister:  context.Core.Services("").Controller().Lister(),
		serviceLister:     context.Rio.Services("").Controller().Lister(),
		endpointLister:    context.Core.Endpoints("").Controller().Lister(),
	}

	c.cache = cache.NewSnapshotCache(true, test.Hasher{}, logrus.StandardLogger())
	srv := server.NewServer(c.cache, nil)

	go test.RunManagementServer(ctx, srv, 5678)

	c.serviceController.AddHandler("envoy-controller", c.sync)

	context.Core.Nodes("").AddHandler("envoy-controller-enque", func(key string, obj *v12.Node) error {
		c.enqueue()
		return nil
	})
	context.Core.Services("").AddHandler("envoy-controller-enque", func(key string, obj *v12.Service) error {
		c.enqueue()
		return nil
	})
	context.Core.Endpoints("").AddHandler("envoy-controller-enque", func(key string, obj *v12.Endpoints) error {
		c.enqueue()
		return nil
	})
}

type controller struct {
	nodeLister        v1.NodeLister
	serviceController v1beta1.ServiceController
	k8sServiceLister  v1.ServiceLister
	serviceLister     v1beta1.ServiceLister
	endpointLister    v1.EndpointsLister
	cache             cache.SnapshotCache
}

func (c *controller) sync(key string, service *v1beta1.Service) error {
	if key == all {
		return c.refresh()
	}
	c.enqueue()
	return nil
}

func (c *controller) refresh() error {
	subdomain := settings.ClusterDomain.Get()
	if subdomain == "" {
		logrus.Infof("waiting for cluster domain to be assigned")
		return nil
	}

	version := fmt.Sprintf("v%d", time.Now().UnixNano())
	clusters, endpoints, err := c.clustersAndEndpoints()
	if err != nil {
		return err
	}

	routes, err := c.routes(subdomain)
	if err != nil {
		return err
	}

	listeners, err := c.listeners()
	if err != nil {
		return err
	}

	snapshot := cache.NewSnapshot(version,
		endpoints,
		clusters,
		routes,
		listeners,
	)

	if err := snapshot.Consistent(); err != nil {
		return err
	}

	return c.cache.SetSnapshot("localhost", snapshot)
}

func (c *controller) clustersAndEndpoints() ([]cache.Resource, []cache.Resource, error) {
	var (
		clusters  []cache.Resource
		endpoints []cache.Resource
	)

	services, err := c.k8sServiceLister.List("", rioServiceSelector)
	if err != nil {
		return nil, nil, err
	}

	for _, service := range services {
		//var targetPort int32
		//for _, port := range service.Spec.Ports {
		//	if port.Name != "http" || port.TargetPort.IntVal <= 0 {
		//		continue
		//	}
		//	targetPort = int32(port.TargetPort.IntVal)
		//}
		//
		//if targetPort <= 0 {
		//	continue
		//}

		clusterName := service.Namespace + ":" + service.Name
		cluster := &v2.Cluster{
			Name:           clusterName,
			ConnectTimeout: 5 * time.Second,
			Type:           v2.Cluster_EDS,
			EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{
						Ads: &core.AggregatedConfigSource{},
					},
				},
			},
		}

		clusters = append(clusters, cluster)

		endpoint, err := c.endpointLister.Get(service.Namespace, service.Name)
		if errors.IsNotFound(err) {
			continue
		}

		for _, subset := range endpoint.Subsets {
			found := false
			for _, endpointPort := range subset.Ports {
				if endpointPort.Port == 80 {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			for _, address := range subset.Addresses {
				endpoints = append(endpoints, newEndpoint(clusterName, address.IP, 80))
			}
		}

	}

	return clusters, endpoints, nil
}

func newEndpoint(clusterName, address string, port int32) *v2.ClusterLoadAssignment {
	return &v2.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: []endpoint.LbEndpoint{{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Protocol: core.TCP,
								Address:  address,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: uint32(port),
								},
							},
						},
					},
				},
			}},
		}},
	}

}

func (c *controller) routes(subdomain string) ([]cache.Resource, error) {
	services, err := c.serviceLister.List("", labels.Everything())
	if err != nil {
		return nil, err
	}

	var virtualHosts []route.VirtualHost
	for _, service := range services {
		port := getHTTPPort(service)
		if port <= 0 {
			continue
		}

		weight := 100
		weights := map[string]int{}

		virtualHosts = append(virtualHosts, directVHosts(service.Namespace, service.Name, "latest", subdomain))

		for rev, revSpec := range service.Spec.Revisions {
			virtualHosts = append(virtualHosts, directVHosts(service.Namespace, service.Name, rev, subdomain))

			revWeight := min(weight, revSpec.Weight)
			if revWeight > 0 {
				weight -= revSpec.Weight
				weights[rev] = revWeight
			}
		}

		if weight > 0 {
			weights["latest"] = weight
		}

		virtualHosts = append(virtualHosts, weightedVHosts(service.Namespace, service.Name, subdomain, weights))
	}

	return []cache.Resource{&v2.RouteConfiguration{
		Name:         "default",
		VirtualHosts: virtualHosts,
	}}, nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func stripRandom(namespace string) string {
	i := strings.LastIndex(namespace, "-")
	if i > 0 {
		return namespace[:i]
	}
	return namespace
}

func weightedVHosts(namespace, serviceName, subdomain string, weights map[string]int) route.VirtualHost {
	host := fmt.Sprintf("%s.%s.%s", serviceName, stripRandom(namespace), subdomain)
	routedWeighted := &route.RouteAction_WeightedClusters{
		WeightedClusters: &route.WeightedCluster{},
	}

	for rev, weight := range weights {
		clusterName := namespace + ":" + serviceName
		if rev != "latest" {
			clusterName += "-" + rev
		}

		routedWeighted.WeightedClusters.Clusters = append(routedWeighted.WeightedClusters.Clusters,
			&route.WeightedCluster_ClusterWeight{
				Weight: &google_protobuf1.UInt32Value{
					Value: uint32(weight),
				},
				Name: clusterName,
			})
	}

	return route.VirtualHost{
		Name:    namespace + ":" + serviceName,
		Domains: []string{host},
		Routes: []route.Route{{
			Match: route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: routedWeighted,
				},
			},
		}},
	}
}

func directVHosts(namespace, serviceName, revision, subdomain string) route.VirtualHost {
	clusterName := namespace + ":" + serviceName
	if revision != "latest" {
		clusterName += "-" + revision
	}

	host := fmt.Sprintf("%s.%s.%s.%s", revision, serviceName, stripRandom(namespace), subdomain)
	return route.VirtualHost{
		Name:    namespace + ":" + serviceName + ":" + revision,
		Domains: []string{host},
		Routes: []route.Route{{
			Match: route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}},
	}
}

func getHTTPPort(service *v1beta1.Service) int32 {
	for _, port := range service.Spec.PortBindings {
		if port.Protocol == "http" {
			//TODO hardcoded
			return 80
		}
	}

	return 0
}

// MakeRoute creates an HTTP route that routes to a given cluster.
func MakeRoute(routeName, clusterName string) *v2.RouteConfiguration {
	return &v2.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []route.VirtualHost{{
			Name:    routeName,
			Domains: []string{"*"},
			Routes: []route.Route{{
				Match: route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
					},
				},
			}},
		}},
	}
}

func (c *controller) listeners() ([]cache.Resource, error) {
	return []cache.Resource{
		MakeHTTPListener("default", 80, "default"),
	}, nil
}

func MakeHTTPListener(listenerName string, port uint32, route string) *v2.Listener {
	// data source configuration
	rdsSource := core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &core.ConfigSource_Ads{
		Ads: &core.AggregatedConfigSource{},
	}

	// access log service configuration
	alsConfig := &als.HttpGrpcAccessLogConfig{
		CommonConfig: &als.CommonGrpcAccessLogConfig{
			LogName: "echo",
			GrpcService: &core.GrpcService{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
						ClusterName: "xds_cluster",
					},
				},
			},
		},
	}
	alsConfigPbst, err := util.MessageToStruct(alsConfig)
	if err != nil {
		panic(err)
	}

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: util.Router,
		}},
		AccessLog: []*alf.AccessLog{{
			Name:   util.HTTPGRPCAccessLog,
			Config: alsConfigPbst,
		}},
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	return &v2.Listener{
		Name: listenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.HTTPConnectionManager,
				Config: pbst,
			}},
		}},
	}
}

func (c *controller) enqueue() {
	c.serviceController.Enqueue("", all)
}
