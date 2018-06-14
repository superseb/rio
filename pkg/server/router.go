package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/settings"
	settings2 "github.com/rancher/rio/pkg/settings"
)

func router(api, k3s http.Handler) http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = k3s
	router.PathPrefix("/v1beta1").Handler(api)
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
