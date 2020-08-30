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
		fmt.Sprintf("relation already exists on inventory: '%+2v'", equipmentsRepeateds),
	})
}

func handleBody(body io.ReadCloser) ([]equipment, error) {
	errs := []string{}
	payloadSingle := equipment{}
	payloadList := []equipment{}
	errSing := json.NewDecoder(body).Decode(&payloadSingle)
	errList := json.NewDecoder(body).Decode(&payloadList)

	if errSing != nil && errList != nil {
		fmt.Println("21", errSing.Error())
		fmt.Println("24", errList.Error())
		logrus.Debugf(fmt.Errorf("11%s\n%s", errSing.Error(), errList.Error()).Error())
		return payloadList, fmt.Errorf("%s\n%s", errSing.Error(), errList.Error())
	}

	if errList != nil {
		logrus.Debugf(fmt.Errorf("22%s\n%s", errSing.Error(), errList.Error()).Error())
		payloadList = append(payloadList, payloadSingle)
	}

	for _, payload := range payloadList {
		if payload.Name == "" || payload.Code == "" { // TODO location tbm é not null?
			logrus.Debugf("payload empty: %+2v", payload)
			errs = append(errs, fmt.Sprintf("payload '%+2v' can't be empty or nil", payload))
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return payloadList, nil
}
