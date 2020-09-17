package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/anuchito/hometic/logger"
	"github.com/gorilla/mux"
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
	r.Use(logger.Middleware)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(&JSONResponseWriter{w}, r)
		})
	})

	r.Handle("/pair-device", CustomHandlerFunc(PairDeviceHandler(NewCreatePairDevice(db)))).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

type JSONResponseWriter struct {
	http.ResponseWriter
}

func (w JSONResponseWriter) json() {
	w.Header().Set("content-type", "application/json")
}

func (w *JSONResponseWriter) Write(p []byte) (int, error) {
	w.json()
	return w.ResponseWriter.Write(p)
}

func (w JSONResponseWriter) WriteHeader(statusCode int) {
	w.json()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *JSONResponseWriter) JSON(statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

type CustomResponseWriter interface {
	http.ResponseWriter
	JSON(statusCode int, data interface{})
}

type CustomHandlerFunc func(CustomResponseWriter, *http.Request)

func (f CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(&JSONResponseWriter{w}, r)
}

func (CustomHandlerFunc) JSON(statusCode int, data interface{}) {
}

func PairDeviceHandler(device Device) func(w CustomResponseWriter, r *http.Request) {
	return func(w CustomResponseWriter, r *http.Request) {
		logger.L(r.Context()).Info("pair-device")
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

		w.JSON(http.StatusOK, map[string]interface{}{"status": "active"})
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
