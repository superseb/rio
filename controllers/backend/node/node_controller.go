package node

import (
	"context"
	"time"

	"github.com/rancher/rancher/pkg/controllers/user/approuter"
	"github.com/rancher/rancher/pkg/ticker"
	"github.com/rancher/rio/pkg/settings"
	"github.com/rancher/rio/types"
	"github.com/rancher/types/apis/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	all         = "_all_"
	clusterName = "rio-system"
)

func Register(ctx context.Context, rContext *types.Context) {
	rdnsClient := approuter.NewClient(rContext.Core.Secrets(""),
		rContext.Core.Secrets("").Controller().Lister(),
		clusterName)
	rdnsClient.SetBaseURL("http://api.lb.rancher.cloud/v1")

	nc := &nodeController{
		rdnsClient:     rdnsClient,
		nodeLister:     rContext.Core.Nodes("").Controller().Lister(),
		nodeController: rContext.Core.Nodes("").Controller(),
	}

	nc.nodeController.AddHandler("node-controller", nc.sync)

	go func() {
		nc.renew()
		for range ticker.Context(ctx, 6*time.Hour) {
			nc.renew()
		}
	}()
}

type nodeController struct {
	rdnsClient     *approuter.Client
	nodeLister     v1.NodeLister
	nodeController v1.NodeController
}

func (n *nodeController) sync(key string, node *corev1.Node) error {
	if node != nil {
		n.nodeController.Enqueue("", all)
		return nil
	}

	if key != all {
		return nil
	}

	ips, err := n.collectIPs()
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return nil
	}

	ips = []string{"127.0.0.1"}

	_, fdqn, err := n.rdnsClient.ApplyDomain(ips)
	if err != nil {
		return err
	}
	settings.ClusterDomain.Set(fdqn)

	return nil
}

func (n *nodeController) renew() error {
	if err := n.sync(all, nil); err != nil {
		return err
	}
	_, err := n.rdnsClient.RenewDomain()
	return err
}

func (n *nodeController) collectIPs() ([]string, error) {
	nodes, err := n.nodeLister.List("", labels.Everything())
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, node := range nodes {
		nodeIP := ""
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeExternalIP {
				nodeIP = address.Address
			} else if address.Type == corev1.NodeInternalIP && nodeIP == "" {
				nodeIP = address.Address
			}
		}

		if nodeIP != "" {
			ips = append(ips, nodeIP)
		}
	}

	return ips, nil
}
