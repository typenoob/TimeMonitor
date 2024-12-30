package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllRecord(w http.ResponseWriter, req *http.Request) {
	s := GetInstance()
	json.NewEncoder(w).Encode(s.GetRecords())
}

func Bootstrap() {
	router := mux.NewRouter()
	router.HandleFunc("/records", GetAllRecord).Methods("GET")
	log.Fatal(http.ListenAndServe(":1234", router))
}
