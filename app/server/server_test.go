package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseRoute = "http://localhost:8092/"

func init() {
	go Run("8092")
}

func TestLiveness(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(liveness)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestReadness(t *testing.T) {
	resp, err := http.Get(fmt.Sprint(baseRoute, "/healthy"))
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		respBody := response{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		t.Errorf("%+2v", respBody)
	}
}
