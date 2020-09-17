package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDeviceHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{ID: 1234, UserID: 4333})
	req := httptest.NewRequest(http.MethodPost, "/devices", payload)
	rec := httptest.NewRecorder()

	handler := CustomHandlerFunc(CreatePairHandler(CreatePairDeviceFunc(func(p Pair) error {
		return nil
	})))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("it should be okay but it not!!!!")
	}

	expected, _ := json.Marshal(map[string]interface{}{"status": "active"})
	if rec.Body.String() != fmt.Sprintf("%s\n", expected) {
		t.Errorf("it should be active you know!!! expected:%q\nbut got:%q\n", rec.Body.String(), expected)
	}
}
