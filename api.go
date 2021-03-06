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
	"flag"
)

type ApiRequestBody struct {
	ImageURL, ImageBody string
}

type ApiResponseBody struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

func main() {
	networkAddress := flag.String("i", "0.0.0.0", "network address")
	networkPort := flag.String("p", "80", "listening port")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api", ApiHandler).Methods("POST")
	http.Handle("/", router)

	fmt.Println(fmt.Sprintf("Up and running, listening on %s:%s", *networkAddress, *networkPort))

	err := http.ListenAndServe(*networkAddress + ":" + *networkPort, nil)

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
	)

	tesseractCommand := exec.Command("tesseract", "stdin", "stdout")

	if apiRequestBody.ImageURL != "" {
		result = parseFromURL(apiRequestBody.ImageURL, tesseractCommand)

		message = fmt.Sprintf("scan image from %s", apiRequestBody.ImageURL)
	} else if apiRequestBody.ImageBody != "" {
		result = parseFromBase64(apiRequestBody.ImageBody, tesseractCommand)

		message = "scan image from base64 body"
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(responseWriter)
	encoder.Encode(ApiResponseBody{
		message,
		result,
	})
}

func parseFromBase64(base64Image string, tesseractCommand *exec.Cmd) string {
	var outputBuffer bytes.Buffer

	imageBody, _ := base64.StdEncoding.DecodeString(base64Image)

	tesseractCommand.Stdin = io.Reader(strings.NewReader(string(imageBody)))
	tesseractCommand.Stdout = &outputBuffer

	tesseractCommand.Start()
	tesseractCommand.Wait()

	return strings.Trim(outputBuffer.String(), "\n")
}

func parseFromURL(imageURL string, tesseractCommand *exec.Cmd) string {
	var outputBuffer bytes.Buffer

	curlCommand := exec.Command("curl", "-s", imageURL)

	pipeReader, pipeWriter := io.Pipe()
	curlCommand.Stdout = pipeWriter

	tesseractCommand.Stdin = pipeReader
	tesseractCommand.Stdout = &outputBuffer

	curlCommand.Start()
	tesseractCommand.Start()
	curlCommand.Wait()
	pipeWriter.Close()
	tesseractCommand.Wait()

	return strings.Trim(outputBuffer.String(), "\n")
}
