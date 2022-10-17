package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// ResponseType represents a type of response
type ResponseType string

// Response types
const (
	Error ResponseType = "error message"
	Info  ResponseType = "info message"
)

type apiResponse struct {
	Code     int          `json:"code"`
	RespType ResponseType `json:"type"`
	Message  string       `json:"message"`
}

// Pet represents a pet
type Pet struct {
	Name    string `json:"name"`
	PetType string `json:"pet_type"`
}

var petDB = []Pet{
	{
		Name:    "Bobby",
		PetType: "cat",
	},
	{
		Name:    "Ralph",
		PetType: "dog",
	},
	{
		Name:    "Shirley",
		PetType: "Armadillo",
	},
}
var dbLock sync.RWMutex

func respondSuccessMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse{
		Code:     http.StatusOK,
		RespType: Info,
		Message:  message,
	})
}

func respondErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(apiResponse{
		Code:     statusCode,
		RespType: Error,
		Message:  message,
	})
}

func respondObject(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object)
}

func petHandler(w http.ResponseWriter, r *http.Request) {
	type petResponse struct {
		ID int `json:"id"`
		Pet
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/pet/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondErrorMessage(w, http.StatusBadRequest, "invalid id")
		return
	}
	switch r.Method {
	case "PUT":
		if r.Header.Get("X-API-KEY") == "" {
			respondErrorMessage(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		dbLock.Lock()
		defer dbLock.Unlock()
		if id < 0 || id >= len(petDB) {
			respondErrorMessage(w, http.StatusNotFound, "pet not found")
			return
		}
		newPet := &Pet{}
		if err := json.NewDecoder(r.Body).Decode(newPet); err != nil {
			respondErrorMessage(w, http.StatusBadRequest, "pet object")
			return
		}
		if newPet.Name != "" {
			petDB[id].Name = newPet.Name
		}
		if newPet.PetType != "" {
			petDB[id].PetType = newPet.PetType
		}
		respondSuccessMessage(w, "pet updated")
		return
	case "GET":
		dbLock.RLock()
		defer dbLock.RUnlock()
		if id < 0 || id >= len(petDB) {
			respondErrorMessage(w, http.StatusNotFound, "pet not found")
			return
		}
		respondObject(w, petResponse{ID: id, Pet: petDB[id]})
		return
	default:
		respondErrorMessage(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}
