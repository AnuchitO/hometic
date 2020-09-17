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
	w.json()
	return w.ResponseWriter.Write(p)
}

func (w *JSONResponseWriter) WriteHeader(statusCode int) {
	w.json()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *JSONResponseWriter) json() {
	w.Header().Set("content-type", "application/json")
}

func (w *JSONResponseWriter) JSON(statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func main() {
	db, err := sql.Open("sqlite3", "hometic.db")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(logger.LoggerMiddleware)
	r.Handle("/pairs", CreatePairHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	srv := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:2009",
	}

	log.Println("staring..")
	log.Fatal(srv.ListenAndServe())
}

type CustomResponseWriter interface {
	http.ResponseWriter
	JSON(statusCode int, data interface{})
}

type CustomHandlerFunc func(CustomResponseWriter, *http.Request)

func (f CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(&JSONResponseWriter{w}, r)
}

func CreatePairHandler(device Device) CustomHandlerFunc {
	return func(w CustomResponseWriter, r *http.Request) {
		logger.L(r.Context()).Info("pair-device")
		var d Pair
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			w.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = device.Pair(d)
		if err != nil {
			w.JSON(http.StatusInternalServerError, err.Error())
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
		_, err := db.Exec("INSERT INTO pairs VALUES (?,?);", p.ID, p.UserID)
		return err
	}
}
