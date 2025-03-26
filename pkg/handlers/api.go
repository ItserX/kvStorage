package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"kvManager/pkg/storage"
)

const (
	ErrIncorrectBody   string = "Incorrect body"
	ErrKeyExists       string = "Key already exists"
	ErrInternalServer  string = "Internal server error"
	ErrKeyIsNotAString string = "Key is not a string"
	ErrReadReqBody     string = "Failed to read request body"
)

type DbHandler struct {
	Repo   storage.TarantoolRepo
	Logger *zap.SugaredLogger
}

type RequestData struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ResponseData struct {
	Value interface{} `json:"value"`
}

func (handler *DbHandler) Add(w http.ResponseWriter, r *http.Request) {
	handler.Logger.Infow("Add request started", "method", r.Method, "path", r.URL.Path)
	data, ok := handler.parseReqBody(w, r)
	if !ok {
		return
	}

	handler.Logger.Debugw("Try to add value", "key", data.Key, "value", data.Value)
	err := handler.Repo.AddValue(data.Key, data.Value)
	if err != nil {
		handler.Logger.Warnw("Falied to add value", "key", data.Key, "value", data.Value,
			"error", err.Error(), "http_status", http.StatusConflict)
		http.Error(w, ErrKeyExists, http.StatusConflict)
		return
	}

	handler.Logger.Infow("Value added successfully", "key", data.Key,
		"http_status", http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}

func (handler *DbHandler) Get(w http.ResponseWriter, r *http.Request) {
	handler.Logger.Infow("Get request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	handler.Logger.Debugw("Try to get value", "key", key)
	data, err := handler.Repo.GetValue(key)
	if handler.checkError(w, err) {
		return
	}

	dataValue := data[0].([]interface{})[1]
	if _, ok := dataValue.(map[interface{}]interface{}); ok {
		dataValue, err = handler.convertMap(dataValue.(map[interface{}]interface{}))
		handler.Logger.Debugw("Try to converting map",
			"map", dataValue)
		if err != nil {
			handler.Logger.Errorw("Converting map failed", "data", dataValue,
				"error", err.Error())
			http.Error(w, ErrKeyIsNotAString, http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(ResponseData{dataValue})
	if err != nil {
		handler.Logger.Errorw("Response marshaling failed", "key", key, "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.Logger.Infow("Get value successful", "key", key,
		"response", string(resp), "http_status", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		handler.Logger.Errorw("Internal server error", "error", err.Error())
	}
}

func (handler *DbHandler) Update(w http.ResponseWriter, r *http.Request) {
	handler.Logger.Infow("Update request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	data, ok := handler.parseReqBody(w, r)
	if !ok {
		return
	}

	handler.Logger.Debugw("Try to get value", "key", key)
	err := handler.Repo.UpdateValue(key, data.Value)
	if handler.checkError(w, err) {
		return
	}

	handler.Logger.Infow("Update value successful", "key", key, "http_status", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func (handler *DbHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handler.Logger.Infow("Delete request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	handler.Logger.Debugw("Try to delete value", "key", key)
	err := handler.Repo.DeleteValue(key)
	if handler.checkError(w, err) {
		return
	}

	handler.Logger.Infow("Delete value successful", "key", key, "http_status", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}
