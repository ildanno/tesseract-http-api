package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"io"
	"fmt"
	"os/exec"
	"bytes"
	"strings"
	"log"
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

	err := http.ListenAndServe(":16680", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("API Doc here"))
}

func ApiHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var apiRequestBody ApiRequestBody

	if err := decoder.Decode(&apiRequestBody); err == io.EOF {
		// ok
	} else if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		result, message string
		outputBuffer    bytes.Buffer
	)

	if apiRequestBody.ImageURL != "" {
		curlCommand := exec.Command("curl", "-s", apiRequestBody.ImageURL)
		tesseractCommand := exec.Command("tesseract", "stdin", "stdout")

		pipeReader, pipeWriter := io.Pipe()
		curlCommand.Stdout = pipeWriter
		tesseractCommand.Stdin = pipeReader

		tesseractCommand.Stdout = &outputBuffer

		curlCommand.Start()
		tesseractCommand.Start()
		curlCommand.Wait()
		pipeWriter.Close()
		tesseractCommand.Wait()

		result = strings.Trim(outputBuffer.String(), "\n")

		message = fmt.Sprintf("scan image from %s", apiRequestBody.ImageURL)
	} else if apiRequestBody.ImageBody != "" {
		imageBody, _ := base64.StdEncoding.DecodeString(apiRequestBody.ImageBody)

		tesseractCommand := exec.Command("tesseract", "stdin", "stdout")

		tesseractCommand.Stdin = io.Reader(strings.NewReader(string(imageBody))) //r

		tesseractCommand.Stdout = &outputBuffer

		tesseractCommand.Start()
		tesseractCommand.Wait()

		result = strings.Trim(outputBuffer.String(), "\n")

		message = "scan image from base64 body"
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(responseWriter)
	encoder.Encode(ApiResponseBody{
		message,
		result,
	})
}
