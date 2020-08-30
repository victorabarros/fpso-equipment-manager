package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

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
		logrus.Errorf("bad request: %s", err.Error())
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

	code := strings.ToUpper(payload.Code) // TODO Case insensitive?

	_, ok := db[code]
	if ok {
		logrus.Errorf("'%s' already exists", code)
		rw.WriteHeader(http.StatusConflict)
		json.NewEncoder(rw).Encode(response{
			fmt.Sprintf("code '%s' already exists", code),
		})
		return
	}

	db[code] = make(map[string]equipment)
	rw.WriteHeader(http.StatusCreated)
}
