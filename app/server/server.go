package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type equipment struct {
	name     string
	code     string
	location string
}

var db = map[string][]equipment{}

// Run up server
func Run(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", liveness)
	r.HandleFunc("/healthy", readiness)
	r.HandleFunc("/vessel", insertVessel).Methods(http.MethodPost)

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

type response struct {
	Message string `json:"message"`
}

func insertVessel(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertVessel\" trigged")
	defer fmt.Printf("%+2v\n", db) // TODO remove
	payload := struct {
		Code string `json:"code"`
	}{}

	rw.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		logrus.Errorf("bad request: %ss", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	} else if payload.Code == "" {
		logrus.Errorf("code empty: %+2v", payload)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"code can't be empty or nil",
		})
		return
	}

	code := strings.ToUpper(payload.Code)

	_, ok := db[code]
	if ok {
		logrus.Errorf("%s already exists", code)
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("code '%s' already exists", code),
		})
		return
	}

	db[code] = []equipment{}
	rw.WriteHeader(http.StatusCreated)
}
