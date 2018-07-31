package gateway

import (
	"strconv"
	"sync"
	"time"

	"strings"

	"context"

	"github.com/rancher/norman/types/set"
	"github.com/rancher/rio/pkg/settings"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/networking.istio.io/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	refreshInterval = 5 * time.Minute
)

type GatewayController struct {
	sync.Mutex

	ports                map[string]bool
	gatewayLister        v1alpha3.GatewayLister
	virtualServiceLister v1alpha3.VirtualServiceLister
	gateways             v1alpha3.GatewayInterface
	lastRefresh          time.Time
}

func Register(ctx context.Context, rContext *types.Context) {
	gc := &GatewayController{
		ports:                map[string]bool{},
		gatewayLister:        rContext.Networking.Gateways("").Controller().Lister(),
		virtualServiceLister: rContext.Networking.VirtualServices("").Controller().Lister(),
		gateways:             rContext.Networking.Gateways(""),
	}
	rContext.Networking.VirtualServices("").Controller().AddHandler("gateway-controller", gc.sync)
}

func (g *GatewayController) sync(key string, service *v1alpha3.VirtualService) error {
	if service == nil {
		return nil
	}

	g.Lock()
	if time.Now().Sub(g.lastRefresh) > refreshInterval {
		g.refresh()
	}
	g.Unlock()

	return g.addPorts(getPorts(service)...)
}

func getPorts(service *v1alpha3.VirtualService) []string {
	ports, ok := service.Annotations["rio.cattle.io/ports"]
	if !ok || ports == "" {
		return nil
	}

	return strings.Split(ports, ",")
}

func (g *GatewayController) refresh() error {
	now := time.Now()
	existingPorts := map[string]bool{}
	newPorts := map[string]bool{}

	gw, err := g.gatewayLister.Get(settings.RioSystemNamespace, settings.IstionExternalGateway)
	if err == nil {
		for _, server := range gw.Spec.Servers {
			existingPorts[strconv.FormatUint(uint64(server.Port.Number), 10)] = true
		}
	}

	vss, err := g.virtualServiceLister.List("", labels.Everything())
	if err != nil {
		return err
	}

	for _, vs := range vss {
		newPorts, _ = addPorts(newPorts, getPorts(vs)...)
	}

	toCreate, toDelete, _ := set.Diff(existingPorts, newPorts)
	if len(toCreate) > 0 || len(toDelete) > 0 {
		err = g.createGateway(newPorts)
	}

	if err != nil {
		return err
	}

	g.lastRefresh = now
	g.ports = newPorts
	return nil
}

func (g *GatewayController) addPorts(ports ...string) error {
	g.Lock()
	defer g.Unlock()

	newPorts, add := addPorts(g.ports, ports...)
	if !add {
		return nil
	}

	return g.createGateway(newPorts)
}

func (g *GatewayController) createGateway(newPorts map[string]bool) error {
	spec := v1alpha3.GatewaySpec{
		Selector: map[string]string{
			"gateway": "external",
		},
	}

	for portStr := range newPorts {
		port, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			continue
		}

		spec.Servers = append(spec.Servers, &v1alpha3.Server{
			Hosts: []string{
				"*",
			},
			Port: &v1alpha3.Port{
				Protocol: "http",
				Number:   uint32(port),
			},
		})
	}

	gw, err := g.gatewayLister.Get(settings.RioSystemNamespace, settings.IstionExternalGateway)
	if errors.IsNotFound(err) {
		if len(spec.Servers) == 0 {
			return nil
		}
		_, err := g.gateways.Create(&v1alpha3.Gateway{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Gateway",
				APIVersion: "networking.istio.io/v1alpha3",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: settings.RioSystemNamespace,
				Name:      settings.IstionExternalGateway,
			},
			Spec: spec,
		})
		return err
	} else if err != nil {
		return err
	}

	if len(spec.Servers) == 0 {
		err = g.gateways.DeleteNamespaced(gw.Namespace, gw.Name, nil)
	} else {
		gw.Spec = spec
		_, err = g.gateways.Update(gw)
	}

	if err != nil {
		return err
	}

	g.ports = newPorts

	return err
}

func addPorts(existingPorts map[string]bool, ports ...string) (map[string]bool, bool) {
	newPorts := map[string]bool{}
	add := false

	for _, port := range ports {
		if _, ok := existingPorts[port]; ok {
			continue
		}
		add = true
		newPorts[port] = true
	}

	if !add {
		return nil, false
	}

	for k, v := range existingPorts {
		newPorts[k] = v
	}

	return newPorts, true
}
