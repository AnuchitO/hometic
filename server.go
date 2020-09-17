package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Pair struct {
	ID     int64
	UserID int64
}

func main() {
	fmt.Println("hello Gopher!")

	r := mux.NewRouter()
	r.HandleFunc("/pair-device", PairDeviceHandler).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}

	log.Println("starting...")
	log.Fatal(server.ListenAndServe())
}

func PairDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var p Pair
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	log.Printf("pair: %#v\n", p)

	w.Write([]byte(`{"status":"active"}`))
}
