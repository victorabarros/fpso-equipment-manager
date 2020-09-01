package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/victorabarros/challenge-modec/internal/database"
)

func TestInsertSingleEquipment(t *testing.T) {
	tests := []struct {
		vessel  string
		body    interface{}
		success bool
	}{
		{
			vessel: "MV102",
			body: database.Equipment{
				Code:     "9074R9W1",
				Location: "Japan",
				Name:     "engine",
			},
			success: true,
		},
		{
			vessel: "mv1020",
			body: database.Equipment{
				Code:     "9074R9W13",
				Location: "Japan",
				Name:     "engine",
			},
			success: true,
		},
		{
			vessel: "MV1021",
			body: struct {
				Code int `json:"code"`
			}{
				4,
			},
			success: false,
		},
		{
			vessel:  "MV1022",
			body:    database.Equipment{},
			success: false,
		},
		{
			vessel: "MV1023",
			body: database.Equipment{
				Code:     "9074R9W1",
				Location: "Japan",
				Name:     "engine",
			},
			success: false,
		},
	}

	for _, test := range tests {
		if err := postVessel(test.vessel); err != nil {
			t.Error("fail to post vessel", err)
		}

		if err := postSingleEquipment(test.vessel, test.body); (err != nil) == test.success {
			t.Error("fail to post single equipment", err)
		}
	}
}
func TestInsertSingleEquipmentNotFound(t *testing.T) {
	vessel := "xpto"
	body := database.Equipment{
		Code:     "9074R9W1",
		Location: "Japan",
		Name:     "engine",
	}

	if err := postSingleEquipment(vessel, body); err == nil {
		t.Error("fail to post single equipment", err)
	}
}

func postSingleEquipment(vessel string, equipt interface{}) error {
	decode, err := json.Marshal(equipt)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprint(baseRoute, "vessel/", vessel, "/equipment")

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(decode))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		respBody := response{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		fmt.Printf("code %d : %+2v\n", resp.StatusCode, respBody)
		return fmt.Errorf("status code '%d': %+2v", resp.StatusCode, respBody)
	}
	return nil
}

func TestInsertEquipmentList(t *testing.T) {
	vessel := "MV103"
	err := postVessel(vessel)
	if err != nil {
		t.Error(err)
	}

	body := []database.Equipment{
		{
			Code:     "5310B9D7",
			Location: "Brazil",
			Name:     "compressor",
		},
		{
			Code:     "1408R2T8",
			Location: "Germany",
			Name:     "pump",
		},
	}
	err = postEquipmentList(vessel, body)
	if err != nil {
		t.Error(err)
	}
}

func postEquipmentList(vessel string, equipts []database.Equipment) error {
	postVessel(vessel)

	body := equipts

	decode, err := json.Marshal(body)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprint(baseRoute, "vessel/", vessel, "/equipments")

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(decode))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func TestFetchEquipment(t *testing.T) {
	vessel := "MV104"
	err := postVessel(vessel)
	if err != nil {
		t.Error(err)
	}

	reqBody := []database.Equipment{
		{
			Code:     "9873B3R7",
			Location: "USA",
			Name:     "tree",
		},
		{
			Code:     "1119T1T5",
			Location: "Italy",
			Name:     "boiler",
		},
	}

	err = postEquipmentList(vessel, reqBody)
	if err != nil {
		t.Error(err)
	}

	endpoint := fmt.Sprint(baseRoute, "vessel/", vessel)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		t.Errorf("%+2v", json.NewDecoder(resp.Body))
	}

	respBody := []database.Equipment{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if len(reqBody) != len(respBody) {
		t.Errorf("body must length %d: %+2v", len(reqBody), respBody)
	}
}

func TestInactiveEquipment(t *testing.T) {
	vessel := "MV105"
	err := postVessel(vessel)
	if err != nil {
		t.Error(err)
	}

	reqBody := database.Equipment{
		Code:     "4319Q1T0",
		Location: "Mexico",
		Name:     "side door",
	}

	err = postSingleEquipment(vessel, reqBody)
	if err != nil {
		t.Error(err)
	}

	endpoint := fmt.Sprint(baseRoute, "equipment/", reqBody.Code)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println(resp.StatusCode)
		respBody := response{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		t.Errorf("%+2v", respBody)
	}

	endpoint = fmt.Sprint(baseRoute, "vessel/", vessel)
	resp, err = http.Get(endpoint)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		t.Errorf("%+2v", json.NewDecoder(resp.Body))
	}

	respBody := []database.Equipment{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if len(respBody) != 0 {
		t.Errorf("body must length 0: %+2v", respBody)
	}
}
