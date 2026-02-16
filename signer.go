package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"
)

// GenerateSignatureForTransactions generates a SNAP signature for transactions using RSA-SHA256
func GenerateSignatureForTransactions(req SignatureRequest) (*SignatureResponse, error) {
	var encodedRequest string

	if req.Body != nil {
		var bodyBytes []byte

		// Check if body is already a string (raw body from Postman)
		if bodyStr, ok := req.Body.(string); ok {
			// Body is already a string, use it directly
			bodyBytes = []byte(bodyStr)
		} else {
			// Body is a map/object, marshal it
			var err error
			bodyBytes, err = json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
		}

		if len(bodyBytes) > 0 {
			// body should be added to signature as Lowercase(HexEncode(SHA-256(minify(RequestBody))))
			var minifiedRequest bytes.Buffer
			err := json.Compact(&minifiedRequest, bodyBytes)
			if err != nil {
				return nil, err
			}
			encodedRequest = GenerateHexEncodedSHA256(minifiedRequest.String())
		}
	}

	if req.Timestamp == "" {
		req.Timestamp = GenerateCurrentTimestamp()
	}

	// decode the key using the appropriate PEM format
	decodedRSAPrivateKey, err := GetPKCS8RSAPrivateKey(context.Background(), req.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Format -> SHA256withRSA (HTTPMethod + ":" + RelativeUrl + ":" + Lowercase(HexEncode(SHA-256(minify(RequestBody)))) + ":" + TimeStamp)
	signatureMessage := strings.Join([]string{req.Method, req.URL, strings.ToLower(encodedRequest), req.Timestamp}, ":")
	signature, err := GenerateBase64EncodedSHA256withRSA(context.Background(), []byte(signatureMessage), decodedRSAPrivateKey)
	if err != nil {
		return nil, err
	}

	return &SignatureResponse{
		Signature:     signature,
		Timestamp:     req.Timestamp,
		StringToSign: signatureMessage,
		Headers: map[string]string{
			"X-TIMESTAMP": req.Timestamp,
			"X-SIGNATURE": signature,
		},
	}, nil
}

func GenerateCurrentTimestamp() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02T15:04:05.000-07:00")
}

// GenerateSignatureForTransactionsHMAC512 generates a SNAP signature for transactions using HMAC-SHA512
func GenerateSignatureForTransactionsHMAC512(req SignatureRequest) (*SignatureResponse, error) {
	var encodedRequest string

	if req.Body != nil {
		var bodyBytes []byte

		// Check if body is already a string (raw body from Postman)
		if bodyStr, ok := req.Body.(string); ok {
			// Body is already a string, use it directly
			bodyBytes = []byte(bodyStr)
		} else {
			// Body is a map/object, marshal it
			var err error
			bodyBytes, err = json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
		}

		if len(bodyBytes) > 0 {
			// body should be added to signature as Lowercase(HexEncode(SHA-256(minify(RequestBody))))
			var minifiedRequest bytes.Buffer
			err := json.Compact(&minifiedRequest, bodyBytes)
			if err != nil {
				return nil, err
			}
			encodedRequest = GenerateHexEncodedSHA256(minifiedRequest.String())
		}
	}

	if req.Timestamp == "" {
		req.Timestamp = GenerateCurrentTimestamp()
	}

	// Format -> HMAC-SHA512 (HTTPMethod + ":" + RelativeUrl + ":" + AccessToken + ":" + Lowercase(HexEncode(SHA-256(minify(RequestBody))))+ ":" + TimeStamp)
	signatureMessage := strings.Join([]string{req.Method, req.URL, req.AccessToken, strings.ToLower(encodedRequest), req.Timestamp}, ":")
	signature, err := GenerateBase64EncodedHMAC512WithSecretKey(req.SecretKey, []byte(signatureMessage))
	if err != nil {
		return nil, err
	}

	return &SignatureResponse{
		Signature:    signature,
		Timestamp:    req.Timestamp,
		StringToSign: signatureMessage,
		Headers: map[string]string{
			"X-TIMESTAMP": req.Timestamp,
			"X-SIGNATURE": signature,
		},
	}, nil
}

// GenerateSignatureForToken generates a SNAP signature for token requests using RSA-SHA256
func GenerateSignatureForToken(req SignatureRequest) (*SignatureResponse, error) {
	if req.Timestamp == "" {
		req.Timestamp = GenerateCurrentTimestamp()
	}

	// decode the key using the appropriate PEM format
	decodedRSAPrivateKey, err := GetPKCS8RSAPrivateKey(context.Background(), req.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Format -> SHA256withRSA (ClientID + "|" + TimeStamp)
	signatureMessage := strings.Join([]string{req.ClientID, req.Timestamp}, "|")
	signature, err := GenerateBase64EncodedSHA256withRSA(context.Background(), []byte(signatureMessage), decodedRSAPrivateKey)
	if err != nil {
		return nil, err
	}

	return &SignatureResponse{
		Signature:    signature,
		Timestamp:    req.Timestamp,
		StringToSign: signatureMessage,
		Headers: map[string]string{
			"X-TIMESTAMP": req.Timestamp,
			"X-SIGNATURE": signature,
		},
	}, nil
}
