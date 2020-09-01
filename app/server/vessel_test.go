package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestInsertVessel(t *testing.T) {
	tests := []struct {
		vessel  string
		success bool
	}{
		{
			vessel:  "MV101",
			success: true,
		},
		{
			vessel:  "MV101",
			success: false,
		},
	}

	for _, test := range tests {
		err := postVessel(test.vessel)
		if (err != nil) == test.success {
			t.Error(err)
		}
	}
}

func TestInsertVesselBadRequest(t *testing.T) {
	tests := []struct {
		body interface{}
	}{
		{
			body: struct {
				Code int `json:"code"`
			}{
				4,
			},
		},
		{
			body: struct {
				Code string `json:"code"`
			}{
				"",
			},
		},
	}

	for _, test := range tests {
		decode, err := json.Marshal(test.body)
		if err != nil {
			t.Error(err)
		}

		endpoint := fmt.Sprint(baseRoute, "vessel")
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(decode))
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			respBody := response{}
			json.NewDecoder(resp.Body).Decode(&respBody)
			t.Error("Fail to make bad request")
		}
	}
}

func postVessel(code string) error {
	body := struct {
		Code string `json:"code"`
	}{code}

	decode, err := json.Marshal(body)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprint(baseRoute, "vessel")
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(decode))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		respBody := response{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		return fmt.Errorf("%+2v", respBody)
	}
	return nil
}
