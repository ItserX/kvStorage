package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

type Handler struct {
	Repo storage.KvRepository
}

func (handler *Handler) Add(w http.ResponseWriter, r *http.Request) {
	log.Logger.Infow("Add request started", "method", r.Method, "path", r.URL.Path)
	data, ok := handler.parseReqBody(w, r)
	if !ok {
		return
	}

	log.Logger.Debugw("Try to add value", "key", data.Key, "value", data.Value)
	err := handler.Repo.AddValue(data.Key, data.Value)
	if err != nil {
		log.Logger.Warnw("Falied to add value", "key", data.Key, "value", data.Value,
			"error", err.Error(), "http_status", http.StatusConflict)
		http.Error(w, ErrKeyExists, http.StatusConflict)
		return
	}

	log.Logger.Infow("Value added successfully", "key", data.Key,
		"http_status", http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}

func (handler *Handler) Get(w http.ResponseWriter, r *http.Request) {
	log.Logger.Infow("Get request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	log.Logger.Debugw("Try to get value", "key", key)
	data, err := handler.Repo.GetValue(key)
	if handler.checkError(w, err) {
		return
	}

	dataValue := data[0].([]any)[1]
	if _, ok := dataValue.(map[any]any); ok {
		dataValue, err = handler.convertMap(dataValue.(map[any]any))
		log.Logger.Debugw("Try to converting map",
			"map", dataValue)
		if err != nil {
			log.Logger.Errorw("Converting map failed", "data", dataValue,
				"error", err.Error())
			http.Error(w, ErrKeyIsNotAString, http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(ResponseData{dataValue})
	if err != nil {
		log.Logger.Errorw("Response marshaling failed", "key", key, "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Logger.Infow("Get value successful", "key", key,
		"response", string(resp), "http_status", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Logger.Errorw("Internal server error", "error", err.Error())
	}
}

func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	log.Logger.Infow("Update request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	data, ok := handler.parseReqBody(w, r)
	if !ok {
		return
	}

	log.Logger.Debugw("Try to get value", "key", key)
	err := handler.Repo.UpdateValue(key, data.Value)
	if handler.checkError(w, err) {
		return
	}

	log.Logger.Infow("Update value successful", "key", key, "http_status", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Logger.Infow("Delete request started", "method", r.Method, "path", r.URL.Path)
	routeVars := mux.Vars(r)
	key := routeVars["id"]

	log.Logger.Debugw("Try to delete value", "key", key)
	err := handler.Repo.DeleteValue(key)
	if handler.checkError(w, err) {
		return
	}

	log.Logger.Infow("Delete value successful", "key", key, "http_status", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}
