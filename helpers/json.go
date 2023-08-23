package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJsonError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if e != nil {
		log.Fatal("json encoding failed")
	}
}

func WriteJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ReadJson(r *http.Request, data any) {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(ServerError{Status: 400, Message: "invalid json request body"})
	}
}

// func ReadJson(r *http.Request, data any) error {
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
