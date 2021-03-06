package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/anuchito/hometic/logger"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

func main() {
	fmt.Println("hello Gopher!")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(logger.Middleware)

	r.Handle("/pair-device", PairDeviceHandler(NewCreatePairDevice(db))).Methods(http.MethodPost)

	addr := fmt.Sprintf("%s:%s", host(os.Getenv("HOST")), port(os.Getenv("PORT")))
	fmt.Println("addr:", addr)
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

func host(h string) string {
	if h == "" {
		return "127.0.0.1"
	}

	return h
}

func port(p string) string {
	if p == "" {
		return "2009"
	}

	return p
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
	JSON(statusCode int, data interface{})
}

type CustomHandlerFunc func(CustomResponseWriter, *http.Request)

func (f CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(&JSONResponseWriter{w}, r)
}

func PairDeviceHandler(device Device) CustomHandlerFunc {
	return func(w CustomResponseWriter, r *http.Request) {
		logger.L(r.Context()).Info("pair-device")
		var p Pair
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer r.Body.Close()

		log.Printf("pair: %#v\n", p)
		if err := device.Pair(p); err != nil {
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
		_, err := db.Exec("INSERT INTO pairs VALUES ($1,$2);", p.DeviceID, p.UserID)
		return err
	}
}
