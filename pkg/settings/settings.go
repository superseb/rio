package settings

import "github.com/rancher/rancher/pkg/settings"

var (
	ClusterDomain = settings.NewSetting("cluster-domain", "")
)
