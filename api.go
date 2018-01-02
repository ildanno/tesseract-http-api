package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"os/exec"
	"bytes"
	"strings"
)

type ApiRequestBody struct {
	ImageURL, ImageBody string
}

type ApiResponseBody struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api", ApiHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println("Up and running, listening on 0.0.0.0:16680")

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

	c1 := exec.Command("curl", "-s", apiRequestBody.ImageURL)
	c2 := exec.Command("tesseract", "stdin", "stdout")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()

	encoder := json.NewEncoder(writer)
	encoder.Encode(ApiResponseBody{
		fmt.Sprintf("execution info for %s", apiRequestBody.ImageURL),
		strings.Trim(b2.String(), "\n"),
	})
}
