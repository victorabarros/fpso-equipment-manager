package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/victorabarros/fpso-equipment-manager/internal/database"
)

func insertSingleEquipment(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertEquipment\" trigged")
	rw.Header().Set("Content-Type", "application/json")

	payload := database.Equipment{}
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		logrus.Errorf("bad request: %ss", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	} else if payload.Name == "" || payload.Code == "" {
		logrus.Errorf("payload empty: %+2v", payload)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"name or code can't be empty or nil",
		})
		return
	}

	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])
	inventory, ok := db[vessel]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}

	payload.Code = strings.ToUpper(payload.Code)

	vesselE, ok := equipmentSet[payload.Code]
	if ok {
		logrus.Errorf("'%s' already exists", payload.Code)
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("equipment already exists on inventory: '%s'", vesselE),
		})
		return
	}

	equipmentSet[payload.Code] = vessel

	payload.Status = "active"

	inventory[payload.Code] = payload
	rw.WriteHeader(http.StatusCreated)
}

func insertEquipmentList(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertEquipment\" trigged")
	rw.Header().Set("Content-Type", "application/json")

	errs := []string{}
	equips := []database.Equipment{}
	err := json.NewDecoder(req.Body).Decode(&equips)
	if err != nil {
		logrus.Debug(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	}

	if len(equips) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: "payload empty",
		})
		return
	}

	for _, equip := range equips {
		if equip.Name == "" || equip.Code == "" {
			logrus.Debugf("equip empty: %+2v", equip)
			errs = append(errs,
				fmt.Sprintf("fields from '%+2v' can't be null or empty", equip))
		}
	}

	if len(errs) > 0 {
		logrus.Debugf("5%+2v\n", errs)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: strings.Join(errs, "\n"),
		})
		return
	}

	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])

	_, ok := db[vessel]
	if !ok {
		logrus.Debugf("vessel '%s' doesn't exists", vessel)
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}

	for _, equip := range equips {
		equip.Code = strings.ToUpper(equip.Code)

		vesselE, ok := equipmentSet[equip.Code]
		if ok {
			errs = append(errs,
				fmt.Sprintf("'%s' already registred to vessel '%s'",
					equip.Code, vesselE))
		}
	}

	if len(errs) > 0 {
		logrus.Debugf("5%+2v\n", errs)
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			Message: strings.Join(errs, "\n"),
		})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	for _, equip := range equips {
		equip.Status = "active"

		equipmentSet[equip.Code] = vessel
		inventory := db[vessel]
		inventory[equip.Code] = equip
	}
}

func fetchEquipments(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"fetchEquipments\" trigged")
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])

	inventory, ok := db[vessel]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}

	resp := []database.Equipment{}
	for _, equip := range inventory {
		if equip.Status == "inactive" {
			continue
		}
		resp = append(resp, equip)
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

func inactiveEquipment(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"patchStatus\" trigged")
	rw.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	equipment := strings.ToUpper(params["equipmentCode"])

	vessel, ok := equipmentSet[equipment]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("equipment '%s' isn't active", equipment),
		})
		return
	}

	inventory := db[vessel]
	data := inventory[equipment]
	data.Status = "inactive"
	inventory[equipment] = data
	delete(equipmentSet, equipment)

	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(response{
		fmt.Sprintf("equipment '%s' inactivated", equipment)})
}
