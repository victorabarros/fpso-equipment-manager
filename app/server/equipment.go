package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	equipments, err := handleBody(req.Body)
	if err != nil {
		logrus.Debugf("bad request: %ss", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	}

	equipmentsUtils := []equipment{}
	equipmentsRepeateds := map[string]string{}

	for _, payload := range equipments {
		payload.Code = strings.ToUpper(payload.Code)

		vesselE, ok := equipmentSet[payload.Code] // TODO improve name
		if ok {
			logrus.Debugf("'%s' already exists", payload.Code)
			equipmentsRepeateds[payload.Code] = vesselE
		} else {
			equipmentsUtils = append(equipmentsUtils, payload)
		}
	}

	for _, payload := range equipmentsUtils {
		payload.Status = "active"

		equipmentSet[payload.Code] = vessel
		inventory[payload.Code] = payload
	}

	if len(equipmentsRepeateds) == 0 {
		rw.WriteHeader(http.StatusCreated)
		return
	}

	if len(equipmentsRepeateds) == len(equipments) {
		rw.WriteHeader(http.StatusConflict)
	} else {
		rw.WriteHeader(http.StatusPartialContent)
	}

	json.NewEncoder(rw).Encode(response{
		fmt.Sprintf("relation already exists on inventory: '%+2v'",
			equipmentsRepeateds),
	})
}

func handleBody(body io.ReadCloser) ([]equipment, error) {
	errs := []string{}
	payloadSingle := equipment{}
	payloadList := []equipment{}
	errList := json.NewDecoder(body).Decode(&payloadList)
	errSing := json.NewDecoder(body).Decode(&payloadSingle)

	if errSing != nil && errList != nil {
		logrus.Debug(errList.Error(), payloadList)
		logrus.Debug(errSing.Error(), payloadSingle)
		return payloadList, fmt.Errorf("%s\n%s", errSing.Error(), errList.Error())
	}

	if errList != nil {
		logrus.Debugf(errList.Error())
		payloadList = append(payloadList, payloadSingle)
	}

	resp := []equipment{}
	for _, payload := range payloadList {
		if payload.Name == "" || payload.Code == "" { // TODO location tbm é not null?
			logrus.Debugf("payload empty: %+2v", payload)
			errs = append(errs, fmt.Sprintf("payload '%+2v' can't be empty or nil", payload))
		} else {
			resp = append(resp, payload)
		}
	}

	if len(errs) > 0 {
		logrus.Debugf("%+2v\n", errs)
		return resp, errors.New(strings.Join(errs, "\n"))
	}

	return resp, nil
}

func fetchEquipments(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"fetchEquipments\" trigged")
	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])
	rw.Header().Set("Content-Type", "application/json")

	inventory, ok := db[vessel]
	if !ok {
		rw.WriteHeader(http.StatusNotFound) // TODO verificar se este é o status code correto
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}
	resp := []equipment{}
	for _, equip := range inventory {
		resp = append(resp, equip)
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

func patchStatus(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"patchStatus\" trigged")
	payload := struct {
		Status string `json:"status"`
	}{}

	rw.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		logrus.Debugf("bad request: %s", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	}

	payload.Status = strings.ToUpper(payload.Status)
	if payload.Status == "" || (payload.Status != "ACTIVE" && payload.Status != "INACTIVE") {
		logrus.Debugf("status invalid: %+2v", payload)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"status invalid",
		})
		return
	}

	params := mux.Vars(req)
	equipment := strings.ToUpper(params["equipmentCode"])

	vessel, ok := equipmentSet[equipment]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("equipment %s not registred", equipment),
		})
		return
	}

	inventory := db[vessel]
	data := inventory[equipment]
	data.Status = payload.Status

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response{
		fmt.Sprintf("status from equipment '%s' updated with success", equipment)})
}
