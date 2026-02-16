package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// respondWithError writes an error response with the given status code and message
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respondWithJSON writes a JSON response with the given status code and data
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func HandleGenerateSignature(w http.ResponseWriter, r *http.Request) {
	var req SignatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Load default private key if not provided (for RSA types)
	if req.PrivateKey == "" {
		defaultKey, err := LoadDefaultPrivateKey()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error loading default private key: %v", err))
			return
		}
		req.PrivateKey = defaultKey
	}

	if err := req.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var resp *SignatureResponse
	var err error

	// Route based on signature type
	switch req.SignatureRequestType {
	case SignatureTypeTransactionsRSA:
		resp, err = GenerateSignatureForTransactions(req)
	case SignatureTypeTransactionsHMAC:
		resp, err = GenerateSignatureForTransactionsHMAC512(req)
	case SignatureTypeToken:
		resp, err = GenerateSignatureForToken(req)
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid signature type")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if default private key is available
	keyLoaded := false
	if _, err := LoadDefaultPrivateKey(); err == nil {
		keyLoaded = true
	}

	response := map[string]any{
		"status":    "ok",
		"keyLoaded": keyLoaded,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
