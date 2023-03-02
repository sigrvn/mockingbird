package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type MockResponse struct {
	ResponseStatus int            `json:"respStatus"`
	ResponseDelay  int            `json:"respDelay"`
	Response       map[string]any `json:"response"`
}

type ConfigureRequest struct {
	TargetApi string       `json:"api"`
	Mock      MockResponse `json:"mock"`
}

var (
	api1Response MockResponse
	api2Response MockResponse
)

type APIResponse struct {
	Status  int    `json:"statusCode"`
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, r APIResponse) {
	w.WriteHeader(r.Status)
	json.NewEncoder(w).Encode(r)
}

func api1Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("got /api1 request")
	w.Header().Set("Content-Type", "application/json")

	log.Printf("mocking a delay of %d ms...\n", api1Response.ResponseDelay)
	time.Sleep(time.Millisecond * time.Duration(api1Response.ResponseDelay))
	log.Printf("sending response for api1...\n", api1Response.ResponseDelay)

	if err := json.NewEncoder(w).Encode(api1Response); err != nil {
		sendResponse(w, APIResponse{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error occured while encoding Response for API1: %s\n", err.Error()),
		})
		return
	}
}

func api2Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("got /api2 request")
	w.Header().Set("Content-Type", "application/json")

	log.Printf("mocking a delay of %d ms...\n", api2Response.ResponseDelay)
	time.Sleep(time.Millisecond * time.Duration(api2Response.ResponseDelay))
	log.Printf("sending response for api2...\n", api2Response.ResponseDelay)

	if err := json.NewEncoder(w).Encode(api2Response); err != nil {
		sendResponse(w, APIResponse{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error occured while encoding Response for API2: %s\n", err.Error()),
		})
		return
	}
}

func cfgHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("got /cfg request")
	if r.Method != http.MethodPost {
		sendResponse(w, APIResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: fmt.Sprintf("Error: Method %s not allowed for endpoint /cfg", r.Method),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req ConfigureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, APIResponse{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error occured while decoding ConfigureRequest: %s\n", err.Error()),
		})
		return
	}

	targetApi := strings.ToLower(req.TargetApi)
	switch targetApi {
	case "api1":
		api1Response = req.Mock
	case "api2":
		api2Response = req.Mock
	default:
		sendResponse(w, APIResponse{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Bad target api for /cfg: %s", targetApi),
		})
		return
	}

	sendResponse(w, APIResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("mock response successfully updated for %s", targetApi),
	})
	return
}

func main() {
	http.HandleFunc("/api1", api1Handler)
	http.HandleFunc("/api2", api2Handler)
	http.HandleFunc("/cfg", cfgHandler)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
