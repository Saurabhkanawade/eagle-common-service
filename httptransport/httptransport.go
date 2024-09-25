package httptransport

import (
	"context"
	"encoding/json"
	"net/http"
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(i)
}

func EncodePostResponse(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(i)
}
