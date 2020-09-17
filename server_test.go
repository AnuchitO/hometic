package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePairDevice(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{DeviceID: 1234, UserID: 455})
	req := httptest.NewRequest(http.MethodPost, "/pair-device", payload)
	rec := httptest.NewRecorder()

	handler := CustomHandlerFunc(PairDeviceHandler(CreatePairDeviceFunc(func(p Pair) error {
		return nil
	})))

	handler.ServeHTTP(rec, req)

	if http.StatusOK != rec.Code {
		t.Error("expect 200 OK but got ", rec.Code)
	}

	expected, _ := json.Marshal(map[string]interface{}{"status": "active"})
	if rec.Body.String() != fmt.Sprintf("%s\n", expected) {
		t.Errorf("it should be active you know!!! expected:%q\nbut got:%q\n", rec.Body.String(), expected)
	}
}
