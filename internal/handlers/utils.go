package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

type RequestData struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type ResponseData struct {
	Value any `json:"value"`
}

func (handler *Handler) checkError(w http.ResponseWriter, err error) bool {
	if err != nil && errors.Is(err, storage.ErrKeyNotFound) {
		http.Error(w, storage.ErrKeyNotFound.Error(), http.StatusNotFound)
		log.Logger.Warnw("Key not found error",
			"http_status", http.StatusNotFound)
		return true
	}
	if err != nil {
		log.Logger.Errorw("Internal server error",
			"http_status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func (handler *Handler) convertMap(oldMap map[any]any) (map[string]any, error) {
	newMap := make(map[string]any)
	log.Logger.Debugw("Starting map conversion", "map_size", len(oldMap))
	for key, val := range oldMap {
		strKey, ok := key.(string)
		if !ok {
			log.Logger.Errorw("Map conversion failed", "error", ErrKeyIsNotAString)
			return nil, fmt.Errorf("%s", ErrKeyIsNotAString)
		}

		if nestedMap, ok := val.(map[any]any); ok {
			log.Logger.Debugw("Converting nested map", "key", strKey)
			convertedNested, err := handler.convertMap(nestedMap)
			if err != nil {
				log.Logger.Errorw("Nested map conversion failed",
					"key", strKey,
					"error", err)
				return nil, err
			}
			newMap[strKey] = convertedNested
		} else {
			newMap[strKey] = val
		}
	}
	log.Logger.Debugw("Map conversion completed", "converted_size", len(newMap))
	return newMap, nil

}
func (handler *Handler) parseReqBody(w http.ResponseWriter, r *http.Request) (*RequestData, bool) {
	log.Logger.Debugw("Parsing request body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorw("Failed to read request body",
			"error", err,
			"http_status", http.StatusInternalServerError)
		http.Error(w, ErrReadReqBody, http.StatusInternalServerError)
		return nil, false
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Logger.Warnw("Request body is not closed", err)
		}
	}()

	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Logger.Warnw("Failed to unmarshal request body",
			"error", err,
			"body", string(body),
			"http_status", http.StatusBadRequest)
		http.Error(w, ErrIncorrectBody, http.StatusBadRequest)
		return nil, false
	}
	log.Logger.Debugw("Request body parsed successfully",
		"data_key", data.Key)
	return &data, true
}
