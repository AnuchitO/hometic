package main

import (
	"database/sql"
	"encoding/json"
	"github.com/anuchito/hometic/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "rsc.io/sqlite"
)

type Pair struct {
	ID     int64
	UserID int64
}

type JSONResponseWriter struct {
	http.ResponseWriter
}

func (w *JSONResponseWriter) Write(p []byte) (int, error) {
	w.Header().Set("content-type", "application/json")
	return w.ResponseWriter.Write(p)
}

func main() {
	db, err := sql.Open("sqlite3", "hometic.db")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(logger.LoggerMiddleware)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(&JSONResponseWriter{w}, r)
		})
	})
	r.Handle("/pairs", CreatePairHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	srv := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:2009",
	}

	log.Println("staring..")
	log.Fatal(srv.ListenAndServe())
}

func CreatePairHandler(device Device) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.L(r.Context()).Info("pair-device")
		var d Pair
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		err = device.Pair(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
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
		_, err := db.Exec("INSERT INTO pairs VALUES (?,?);", p.ID, p.UserID)
		return err
	}
}
