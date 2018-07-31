package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"strings"

	"github.com/rancher/norman/api"
	"github.com/rancher/norman/leader"
	"github.com/rancher/norman/signal"
	"github.com/rancher/rancher/k8s"
	"github.com/rancher/rancher/pkg/dynamiclistener"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/rancher/pkg/tls"
	"github.com/rancher/rio/api/setup"
	"github.com/rancher/rio/cli/pkg/resolvehome"
	"github.com/rancher/rio/controllers/backend"
	"github.com/rancher/rio/pkg/data"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/space.cattle.io/v1beta1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	net2 "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/kubernetes/cmd/server"
)

func k3sConfig(dataDir string) (*server.ServerConfig, http.Handler, error) {
	dataDir, err := resolvehome.Resolve(dataDir)
	if err != nil {
		return nil, nil, err
	}

	listenIP := net.ParseIP("127.0.0.1")
	_, clusterIPNet, _ := net.ParseCIDR("10.42.0.0/16")
	_, serviceIPNet, _ := net.ParseCIDR("10.43.0.0/16")

	return &server.ServerConfig{
		PublicIP:       &listenIP,
		PublicPort:     6444,
		PublicHostname: "localhost",
		ListenAddr:     listenIP,
		ListenPort:     6443,
		ClusterIPRange: *clusterIPNet,
		ServiceIPRange: *serviceIPNet,
		DataDir:        dataDir,
	}, newTunnel(), nil
}

func StartServer(ctx context.Context, dataDir string, httpPort, httpsPort int, controllers, inCluster bool) error {
	ctx = signal.SigTermCancelContext(ctx)

	sc, tunnel, err := k3sConfig(dataDir)
	if err != nil {
		return err
	}
	ctx = k8s.SetK3sConfig(ctx, sc)

	embedded, ctx, restConfig, err := k8s.GetConfig(ctx, "auto", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	rContext, err := types.NewContext(*restConfig)
	if err != nil {
		return err
	}
	rContext.Embedded = embedded

	if err := setup.SetupTypes(ctx, rContext); err != nil {
		return err
	}

	apiServer := api.NewAPIServer()
	if err := apiServer.AddSchemas(rContext.Schemas); err != nil {
		return err
	}

	apiRContext, err := types.NewContext(*restConfig)
	if err != nil {
		return err
	}
	apiRContext.Schemas = rContext.Schemas
	apiRContext.Embedded = embedded

	if controllers {
		go leader.RunOrDie(ctx, "rio-controllers", rContext.K8s, func(ctx context.Context) {
			if err := data.AddData(rContext, inCluster); err != nil {
				panic(err)
			}

			if err := backend.Register(ctx, rContext); err != nil {
				panic(err)
			}

			if err := rContext.Start(ctx); err != nil {
				panic(err)
			}

			<-ctx.Done()
		})
	}

	root := router(sc, apiServer, sc.Handler, tunnel)

	if err := startServer(ctx, apiRContext, httpPort, httpsPort, root); err != nil {
		return err
	}

	if err := apiRContext.Start(ctx); err != nil {
		return err
	}

	var (
		clientFile string
		nodeFile   string
	)

	if len(sc.ClientToken) > 0 {
		p := filepath.Join(sc.DataDir, "client-token")
		if err := writeToken(sc.ClientToken, p); err != nil {
			return err
		}
		logrus.Infof("Client token is available at %s", p)
		clientFile = p
	}

	if len(sc.NodeToken) > 0 {
		p := filepath.Join(sc.DataDir, "node-token")
		if err := writeToken(sc.NodeToken, p); err != nil {
			return err
		}
		logrus.Infof("Node token is available at %s", p)
		nodeFile = p
	}

	if len(clientFile) > 0 {
		printToken(httpsPort, "To use CLI:", clientFile, "login")
	}

	if len(nodeFile) > 0 {
		printToken(httpsPort, "To join node to cluster:", nodeFile, "agent")
	}

	<-ctx.Done()
	return nil
}

func printToken(httpsPort int, prefix, file, cmd string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Error(err)
		return err
	}

	token := strings.TrimSpace(string(content))
	ip, err := net2.ChooseHostInterface()
	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Infof("%s rio %s -s https://%s:%d -t %s", prefix, cmd, ip.String(), httpsPort, token)
	return nil
}

func writeToken(token, file string) error {
	if len(token) == 0 {
		return nil
	}

	prefix := "R10"
	certs := settings.CACerts.Get()
	if len(certs) > 0 {
		digest := sha256.Sum256([]byte(certs))
		prefix = "R10" + hex.EncodeToString(digest[:]) + "::"
	}

	return ioutil.WriteFile(file, []byte(prefix+token+"\n"), 0600)
}

func startServer(ctx context.Context, rContext *types.Context, httpPort, httpsPort int, handler http.Handler) error {
	s := &storage{
		listenConfigs:      rContext.Global.ListenConfigs(""),
		listenConfigLister: rContext.Global.ListenConfigs("").Controller().Lister(),
	}
	s2 := &storage2{
		listenConfigs: s.listenConfigs,
	}

	lc, err := tls.ReadTLSConfig(nil)
	if err != nil {
		return err
	}

	if err := tls.SetupListenConfig(s2, false, lc); err != nil {
		return err
	}

	server := dynamiclistener.NewServer(ctx, s, handler, httpPort, httpsPort)
	settings.CACerts.Set(lc.CACerts)
	_, err = server.Enable(lc)
	return err
}

type storage2 struct {
	listenConfigs v1beta1.ListenConfigInterface
}

func (s *storage2) Create(lc *v3.ListenConfig) (*v3.ListenConfig, error) {
	createLC := &v1beta1.ListenConfig{
		ListenConfig: *lc,
	}
	createLC.APIVersion = "space.cattle.io/v1beta1"

	result, err := s.listenConfigs.Create(createLC)
	if err != nil {
		return nil, err
	}
	return &result.ListenConfig, nil
}

func (s *storage2) Get(name string, opts metav1.GetOptions) (*v3.ListenConfig, error) {
	lc, err := s.listenConfigs.Get(name, opts)
	if err != nil {
		return nil, err
	}
	return &lc.ListenConfig, nil
}

func (s *storage2) Update(lc *v3.ListenConfig) (*v3.ListenConfig, error) {
	updateLC := &v1beta1.ListenConfig{
		ListenConfig: *lc,
	}
	updateLC.APIVersion = "space.cattle.io/v1beta1"

	result, err := s.listenConfigs.Update(updateLC)
	if err != nil {
		return nil, err
	}
	return &result.ListenConfig, nil
}

type storage struct {
	listenConfigs      v1beta1.ListenConfigInterface
	listenConfigLister v1beta1.ListenConfigLister
}

func (s *storage) Update(lc *v3.ListenConfig) (*v3.ListenConfig, error) {
	updateLC := &v1beta1.ListenConfig{
		ListenConfig: *lc,
	}
	updateLC.APIVersion = "space.cattle.io/v1beta1"

	updateLC, err := s.listenConfigs.Update(updateLC)
	if err != nil {
		return nil, err
	}
	return &updateLC.ListenConfig, nil
}

func (s *storage) Get(namespace, name string) (*v3.ListenConfig, error) {
	lc, err := s.listenConfigLister.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	return &lc.ListenConfig, nil
}
