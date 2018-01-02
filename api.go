package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
	Result string `json:"result"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api", ApiHandler).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe(":16680", nil)
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("API Doc here"))
}

func ApiHandler(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.Encode(Response{
		"execution info",
		"parsed text",
	})
}
