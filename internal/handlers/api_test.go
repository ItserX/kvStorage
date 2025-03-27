package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"

	"kvManager/internal/handlers"
	"kvManager/internal/mocks"
	"kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

type Case struct {
	Path           string
	Body           string
	Method         string
	ExpectedStatus int
}

func TestAPIHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	err := log.SetupLogger()
	if err != nil {
		t.Errorf("failed to initialize logger: %v", err)
		return
	}
	mockRepo := mocks.NewMockKvRepository(ctrl)

	handler := handlers.Handler{Repo: mockRepo}

	router := mux.NewRouter()
	router.HandleFunc("/kv", handler.Add).Methods("POST")
	router.HandleFunc("/kv/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/kv/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/kv/{id}", handler.Delete).Methods("DELETE")

	testCases := []struct {
		method         string
		path           string
		body           string
		mockSetup      func()
		expectedStatus int
	}{
		{
			method: "POST",
			path:   "/kv",
			body:   `{"key":"test1", "value":{"k1":123, "k2":true}}`,
			mockSetup: func() {
				mockRepo.EXPECT().
					AddValue("test1", map[string]any{"k1": float64(123), "k2": true}).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			method: "POST",
			path:   "/kv",
			body:   `{"key":"test2", "value":{"k1":123}}`,
			mockSetup: func() {
				mockRepo.EXPECT().
					AddValue("test2", gomock.Any()).
					Return(errors.New(handlers.ErrKeyExists)).
					Times(1)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			method: "GET",
			path:   "/kv/test1",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetValue("test1").
					Return([]any{[]any{"test1", map[string]any{"k1": 123}}}, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			method: "GET",
			path:   "/kv/non",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetValue("non").
					Return(nil, storage.ErrKeyNotFound).
					Times(1)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			method: "PUT",
			path:   "/kv/test1",
			body:   `{"value":"new_value"}`,
			mockSetup: func() {
				mockRepo.EXPECT().
					UpdateValue("test1", "new_value").
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			method: "DELETE",
			path:   "/kv/test1",
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteValue("test1").
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			method: "POST",
			path:   "/kv",
			body:   `{"key": "	test4", {invalid}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			if tc.mockSetup != nil {
				tc.mockSetup()
			}

			var req *http.Request
			var err error

			if tc.body != "" {
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}
