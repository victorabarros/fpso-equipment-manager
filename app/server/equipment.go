package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func insertEquipment(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertEquipment\" trigged")
	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])

	inventory, ok := db[vessel]
	rw.Header().Set("Content-Type", "application/json")
	if !ok {
		logrus.Debugf("vessel '%s' doesn't exists", vessel)
		rw.WriteHeader(http.StatusNotFound) // TODO verificar se este é o status code correto
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}

	payload := equipment{}
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		logrus.Debugf("bad request: %ss", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	} else if payload.Name == "" || payload.Location == "" || payload.Code == "" {
		logrus.Debugf("payload empty: %+2v", payload)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"payload can't be empty or nil",
		})
		return
	}

	payload.Code = strings.ToUpper(payload.Code)

	vesselE, ok := equipmentSet[payload.Code] // TODO improve name
	if ok {
		logrus.Debugf("'%s' already exists", payload.Code)
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("'%s' already exists on inventory from vessel '%s'", payload.Code, vesselE),
		})
		return
	}

	equipmentSet[payload.Code] = vessel

	payload.Status = true

	inventory[payload.Code] = payload
}
