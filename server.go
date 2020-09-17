package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/anuchito/hometic/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	_ "rsc.io/sqlite"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

func main() {
	fmt.Println("hello Gopher!")
	db, err := sql.Open("sqlite3", "hometic.db")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(logger.LoggerMiddleware)

	r.Handle("/pair-device", PairDeviceHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

func PairDeviceHandler(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Context().Value("logger").(*zap.Logger).Info("pair-device")
		var p Pair
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		defer r.Body.Close()

		log.Printf("pair: %#v\n", p)
		err = device.Pair(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Write([]byte(`{"status":"active"}`))
	}
}

type Device interface {
	Pair(p Pair) error
}

type CreatePairDeviceFunc func(p Pair) error

func (fn CreatePairDeviceFunc) Pair(p Pair) error {
	return fn(p)
}

func NewCreatePairDevice(db *sql.DB) CreatePairDeviceFunc {
	return func(p Pair) error {
		_, err := db.Exec("INSERT INTO pairs VALUES (?,?);", p.DeviceID, p.UserID)
		return err
	}
}
