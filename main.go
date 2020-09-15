package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	w.Write([]byte(`{"status":"active"}`))
}
