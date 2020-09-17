package main

import (
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
	w.Write([]byte(`{"status":"active"}`))
}
