package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
)

type ApiRequestBody struct {
	ImageURL, ImageBody string
}

type ApiResponseBody struct {
	Message string `json:"message"`
	Result  string `json:"result"`
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
	decoder := json.NewDecoder(request.Body)
	var apiRequestBody ApiRequestBody

	if err := decoder.Decode(&apiRequestBody); err == io.EOF {
		// ok
	} else if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(writer)
	encoder.Encode(ApiResponseBody{
		fmt.Sprintf("execution info for %s", apiRequestBody.ImageURL),
		"parsed text",
	})
}
