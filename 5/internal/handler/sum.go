package handler

import (
	"encoding/json"
	"io"
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
	if err != nil && err != io.EOF {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	resp := Response{Sum: req.A + req.B}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
