package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/settings"
	settings2 "github.com/rancher/rio/pkg/settings"
	"k8s.io/kubernetes/cmd/server"
)

func router(serverConfig *server.ServerConfig, api, k3s http.Handler) http.Handler {
	authed := mux.NewRouter()
	authed.Use(authMiddleware(serverConfig))
	authed.NotFoundHandler = k3s
	authed.Path("/v1beta1/connect").Handler(newTunnel())
	authed.PathPrefix("/v1beta1").Handler(api)
	authed.Path("/node.crt").Handler(nodeCrt(serverConfig))
	authed.Path("/node.key").Handler(nodeKey(serverConfig))

	router := mux.NewRouter()
	router.NotFoundHandler = authed
	router.Path("/cacerts").Handler(cacerts())
	router.Path("/domain").Handler(domain())

	return router
}

func cacerts() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("content-type", "text/plain")
		resp.Write([]byte(settings.CACerts.Get()))
	})
}

func domain() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("content-type", "text/plain")
		resp.Write([]byte(settings2.ClusterDomain.Get()))
	})
}

func nodeCrt(server *server.ServerConfig) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, server.NodeCert)
	})
}

func nodeKey(server *server.ServerConfig) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.TLS == nil {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(resp, req, server.NodeKey)
	})
}
