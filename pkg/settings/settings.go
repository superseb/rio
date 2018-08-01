package settings

import "github.com/rancher/rancher/pkg/settings"

const (
	RioSystemNamespace    = "rio-system"
	RioDefaultNamespace   = "rio-defaults"
	IstionConfigMapName   = "mesh"
	IstionConfigMapKey    = "content"
	IstionExternalGateway = "external"
)

var (
	Version        = "dev"
	ClusterDomain  = settings.NewSetting("cluster-domain", "")
	IstioStackName = settings.NewSetting("istio-stack-name", "istio")
	IstioEnabled   = settings.NewSetting("istio", "true")
	RDNSURL        = settings.NewSetting("rdns-url", "http://api.lb.rancher.cloud/v1")
	RioImage       = settings.NewSetting("rio-image", "rancher/rio")
)

func RioFullImage() string {
	return RioImage.Get() + ":" + Version
}
