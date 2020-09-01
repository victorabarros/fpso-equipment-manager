package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestInsertVessel(t *testing.T) {
	err := postVessel("MV101")
	if err != nil {
		t.Error(err)
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
		return fmt.Errorf("%+2v", json.NewDecoder(resp.Body))
	}
	return nil
}
