package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	if err != nil && storage.ErrKeyNotFound.Error() == err.Error() {
		http.Error(w, storage.ErrKeyNotFound.Error(), http.StatusNotFound)
		handler.Logger.Warnw("Key not found error",
			"http_status", http.StatusNotFound)
		return true
	}
	if err != nil {
		handler.Logger.Errorw("Internal server error",
			"http_status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func (handler *Handler) convertMap(oldMap map[interface{}]interface{}) (map[string]interface{}, error) {
	newMap := make(map[string]interface{})
	handler.Logger.Debugw("Starting map conversion", "map_size", len(oldMap))
	for key, val := range oldMap {
		strKey, ok := key.(string)
		if !ok {
			handler.Logger.Errorw("Map conversion failed", "error", ErrKeyIsNotAString)
			return nil, fmt.Errorf("%s", ErrKeyIsNotAString)
		}

		if nestedMap, ok := val.(map[interface{}]interface{}); ok {
			handler.Logger.Debugw("Converting nested map", "key", strKey)
			convertedNested, err := handler.convertMap(nestedMap)
			if err != nil {
				handler.Logger.Errorw("Nested map conversion failed",
					"key", strKey,
					"error", err)
				return nil, err
			}
			newMap[strKey] = convertedNested
		} else {
			newMap[strKey] = val
		}
	}
	handler.Logger.Debugw("Map conversion completed", "converted_size", len(newMap))
	return newMap, nil

}
func (handler *Handler) parseReqBody(w http.ResponseWriter, r *http.Request) (*RequestData, bool) {
	handler.Logger.Debugw("Parsing request body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.Logger.Errorw("Failed to read request body",
			"error", err,
			"http_status", http.StatusInternalServerError)
		http.Error(w, ErrReadReqBody, http.StatusInternalServerError)
		return nil, false
	}
	defer r.Body.Close()

	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		handler.Logger.Warnw("Failed to unmarshal request body",
			"error", err,
			"body", string(body),
			"http_status", http.StatusBadRequest)
		http.Error(w, ErrIncorrectBody, http.StatusBadRequest)
		return nil, false
	}
	handler.Logger.Debugw("Request body parsed successfully",
		"data_key", data.Key)
	return &data, true
}
