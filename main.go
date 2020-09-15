package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "rsc.io/sqlite"
)

type Pair struct {
	ID     int64
	UserID int64
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/pairs", CreateDevicesHandler).Methods(http.MethodPost)

	srv := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:2009",
	}

	log.Println("staring..")
	log.Fatal(srv.ListenAndServe())
}

func CreateDevicesHandler(w http.ResponseWriter, r *http.Request) {
	var d Pair
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	db, err := sql.Open("sqlite3", "hometic.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO pairs VALUES (?,?);", d.ID, d.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Write([]byte(`{"status":"active"}`))
}
