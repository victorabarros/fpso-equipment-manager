package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type equipment struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

// db[vesselCode][equipmentCode]
var db = map[string]map[string]equipment{}
var equipmentSet = make(map[string]string)

// Run up server
func Run(port string) {
	db["E"] = map[string]equipment{} // TODO remove
	r := mux.NewRouter()
	r.HandleFunc("/healthz", liveness)
	r.HandleFunc("/healthy", readiness)
	r.HandleFunc("/vessel", insertVessel).Methods(http.MethodPost)
	r.HandleFunc("/vessel/{vesselCode}", fetchEquipments).Methods(http.MethodGet)
	r.HandleFunc("/vessel/{vesselCode}/equipment", insertSingleEquipment).Methods(http.MethodPost)
	r.HandleFunc("/vessel/{vesselCode}/equipments", insertEquipmentList).Methods(http.MethodPost)
	r.HandleFunc("/equipment/{equipmentCode}", inactiveEquipment).Methods(http.MethodDelete)

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
	logrus.Debug("route \"healthz\" trigged")

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	json.NewEncoder(rw).Encode(struct {
		ServiceName string
		Version     string
	}{
		"Modec Challenge",
		"v1.0.0",
	})
}

// readiness is k8S readiness probe, returns if pod is read te recieve traffic
func readiness(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"healthy\" trigged")

	// check if all dependencies are alives. (db, services, cache, ...)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	json.NewEncoder(rw).Encode([]struct {
		ServiceName string
		Success     bool
	}{})
}
