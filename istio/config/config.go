package config

//import (
//	"sync"
//
//	"bytes"
//
//	"github.com/golang/protobuf/ptypes/duration"
//	"istio.io/api/mesh/v1alpha1"
//	"istio.io/istio/pilot/pkg/kube/inject"
//	"k8s.io/api/core/v1"
//)
//
//var (
//	templateOnce  = sync.Once{}
//	templateValue string
//)
//
//func MeshConfig() *v1alpha1.MeshConfig {
//	return &v1alpha1.MeshConfig{
//		AuthPolicy: v1alpha1.MeshConfig_MUTUAL_TLS,
//		MtlsExcludedServices: []string{
//			"kubernetes.default.svc.cluster.local",
//		},
//		EnableTracing:     true,
//		MixerCheckServer:  "istio-policy.rio-system:15004",
//		MixerReportServer: "istio-telemetry.rio-system:15004",
//		IngressService:    "rio-ingress",
//		RdsRefreshDelay: &duration.Duration{
//			Seconds: 10,
//		},
//		DefaultConfig: &v1alpha1.ProxyConfig{
//			DiscoveryRefreshDelay: &duration.Duration{
//				Seconds: 10,
//			},
//			ConnectTimeout: &duration.Duration{
//				Seconds: 10,
//			},
//			ConfigPath: "/etc/istio/proxy",
//			BinaryPath: "/usr/local/bin/envoy",
//			DrainDuration: &duration.Duration{
//				Seconds: 45,
//			},
//			ParentShutdownDuration: &duration.Duration{
//				Seconds: 60,
//			},
//			ProxyAdminPort:         15000,
//			ZipkinAddress:          "zipkin.rio-system:9411",
//			StatsdUdpAddress:       "istio-statsd-prom-bridge.rio-system:9125",
//			ControlPlaneAuthPolicy: v1alpha1.AuthenticationPolicy_MUTUAL_TLS,
//			DiscoveryAddress:       "istio-pilot.rio-system:15005",
//		},
//	}
//}
//
//func InjectParams() *inject.Params {
//	debug := true
//	hub := "docker.io/istio"
//	tag := "0.8.latest"
//
//	return &inject.Params{
//		InitImage:           inject.InitImageName(hub, tag, debug),
//		ProxyImage:          inject.ProxyImageName(hub, tag, debug),
//		Verbosity:           2,
//		SidecarProxyUID:     uint64(1337 + 1), // that's right, one better
//		Version:             "",
//		EnableCoreDump:      false,
//		Mesh:                MeshConfig(),
//		ImagePullPolicy:     string(v1.PullIfNotPresent),
//		IncludeIPRanges:     "*",
//		ExcludeIPRanges:     "",
//		IncludeInboundPorts: "*",
//		ExcludeInboundPorts: "",
//		DebugMode:           debug,
//	}
//}
//
//func template() string {
//	templateOnce.Do(func() {
//		var err error
//		templateValue, err = inject.GenerateTemplateFromParams(InjectParams())
//		if err != nil {
//			panic(err)
//		}
//	})
//
//	return templateValue
//}

//func Inject(input []byte) ([]byte, error) {
// 	in := bytes.NewBuffer(input)
//	out := bytes.NewBuffer(make([]byte, 0, len(input)))
//	err := inject.IntoResourceFile(template(), MeshConfig(), in, out)
//	return out.Bytes(), err
//}
