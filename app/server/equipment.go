package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func insertSingleEquipment(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertEquipment\" trigged")
	rw.Header().Set("Content-Type", "application/json")
	defer fmt.Printf("%+2v\n", db) // TODO remove

	payload := equipment{}
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		logrus.Errorf("bad request: %ss", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	} else if payload.Name == "" || payload.Code == "" {
		//TODO documentar que deixou location nullable
		logrus.Errorf("payload empty: %+2v", payload)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"name, location, code or vessel can't be empty or nil",
		})
		return
	}

	params := mux.Vars(req)
	vessel := strings.ToUpper(params["vesselCode"])
	if vessel == "EQUIPMENT" {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			"vessel code empty",
		})
		return
	}

	inventory, ok := db[vessel]
	if !ok {
		logrus.Errorf("vessel '%s' doesn't exists", vessel)
		rw.WriteHeader(http.StatusNotFound) // TODO verificar se este é o status code correto
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("vessel '%s' doesn't exists", vessel),
		})
		return
	}

	payload.Code = strings.ToUpper(payload.Code)

	vesselE, ok := equipmentSet[payload.Code] // TODO improve name
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
}

func insertEquipmentList(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insertEquipment\" trigged")
	rw.Header().Set("Content-Type", "application/json")
	defer fmt.Printf("%+2v\n", db) // TODO remove

	errs := []string{}
	equips := []equipment{}
	err := json.NewDecoder(req.Body).Decode(&equips)

	if err != nil {
		logrus.Debug(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			Message: err.Error(),
		})
		return
	}

	for _, equip := range equips {
		if equip.Name == "" || equip.Code == "" {
			//TODO documentar que deixou location nullable
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
	// TODO validar se vessel é empty

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

		vesselE, ok := equipmentSet[equip.Code] // TODO improve name
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
	defer fmt.Printf("%+2v\n", db) // TODO remove
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

	params := mux.Vars(req)
	equipment := strings.ToUpper(params["equipmentCode"])

	vessel, ok := equipmentSet[equipment]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("equipment '%s' isn't registred", equipment),
		})
		return
	}

	inventory := db[vessel]
	data := inventory[equipment]
	data.Status = "inactive"
	// TODO remover equipmentSet[equipment]

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response{
		fmt.Sprintf("status from equipment '%s' updated with success", equipment)})
}
