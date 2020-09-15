package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	_ "rsc.io/sqlite"
)

type Pair struct {
	ID     int64
	UserID int64
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, _ := zap.NewDevelopment()
		l = l.With(zap.Namespace("hometic"), zap.String("I'm", "gopher"))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", l)))
	})
}

func main() {
	db, err := sql.Open("sqlite3", "hometic.db")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(LoggerMiddleware)
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
		r.Context().Value("logger").(*zap.Logger).Info("pair-device")
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
