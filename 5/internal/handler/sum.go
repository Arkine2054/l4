package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Response struct {
	Sum int `json:"sum"`
}

func SumHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error decoding request: %v", err)
	}

	resp := Response{Sum: req.A + req.B}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
