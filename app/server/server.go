package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Run up server
func Run(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", liveness)
	r.HandleFunc("/healthy", readiness)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	logrus.Debugf("Up apllication at port %s\n", port)
	panic(srv.ListenAndServe())
}

// liveness is k8S liveness probe, returns if pod is alive
// Inspired on: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
func liveness(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"healthz\" started")

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(struct {
		ServiceName string
		Version     string
	}{
		"Modec Challenge",
		"v1.0.0",
	})
}

// readiness is k8S readiness probe, returns if pod is read te recieve traffic
// Inspired on: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
func readiness(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"healthy\" started")

	// check if all dependencies are alives. (db, services, cache, ...)

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode([]struct {
		ServiceName string
		Success     bool
	}{})
}
