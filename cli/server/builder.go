package server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	spaceclient "github.com/rancher/rio/types/client/space/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ContextBuilder struct {
	cfg        *rest.Config
	prefix     string
	serverURL  url.URL
	httpClient *http.Client
}

func (c *ContextBuilder) Domain() (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.url("/domain"), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	domain, err := ioutil.ReadAll(resp.Body)
	return string(domain), err
}

func (c *ContextBuilder) Client(space string) (*client.Client, error) {
	return client.NewClient(&clientbase.ClientOpts{
		URL:        c.url("/v1beta1-rio/spaces/" + space + "/schemas"),
		HTTPClient: c.httpClient,
	})
}

func (c *ContextBuilder) SpaceClient() (*spaceclient.Client, error) {
	return spaceclient.NewClient(&clientbase.ClientOpts{
		URL:        c.url("/v1beta1-rio/schemas"),
		HTTPClient: c.httpClient,
	})
}

func (c *ContextBuilder) url(p string) string {
	newURL := c.serverURL
	newURL.Path = path.Join(c.prefix, p)
	return newURL.String()
}

func NewContextBuilder(config string, k8s bool) (*ContextBuilder, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}

	prefix := ""
	if k8s {
		prefix = "/api/v1/namespaces/rio-system/services/https:rio:https/proxy"
	}

	if strings.Contains(cfg.Host, "/") {
		u, err := url.Parse(cfg.Host)
		if err == nil {
			prefix = path.Join(u.Path, prefix)
		}
	}

	rt, err := rest.TransportFor(cfg)
	if err != nil {
		return nil, err
	}

	if len(prefix) > 1 { // ignore prefix=/
		rt = newCallback(rt, func(req *http.Request) {
			req.Header.Set("X-API-URL-Prefix", prefix)
		})
	}

	url, _, err := rest.DefaultServerURL(cfg.Host, "", schema.GroupVersion{}, true)
	if err != nil {
		return nil, err
	}

	return &ContextBuilder{
		cfg:    cfg,
		prefix: prefix,
		httpClient: &http.Client{
			Transport: rt,
		},
		serverURL: *url,
	}, nil
}
