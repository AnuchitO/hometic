package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPairDevice struct{}

func (mockPairDevice) Pair(p Pair) error  {
	return nil
}

func TestCreateDeviceHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{ID: 1234, UserID: 4333})
	req := httptest.NewRequest(http.MethodPost, "/devices", payload)
	rec := httptest.NewRecorder()

	handler := CreatePairHandler(mockPairDevice{})

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("it should be okay but it not!!!!")
	}

	if rec.Body.String() != `{"status":"active"}` {
		t.Error("it should be active you know!!!")
	}
}
