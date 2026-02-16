package main

// SignatureRequest represents the incoming request to generate a signature
type SignatureRequest struct {
	SignatureRequestType SignatureType `json:"signatureRequestType"`

	// For transactions (both RSA and HMAC)
	Method      string `json:"method,omitempty"`
	URL         string `json:"url,omitempty"`
	Body        any    `json:"body,omitempty"` // Can be map[string]any or string

	// For token generation
	ClientID    string `json:"clientID,omitempty"`

	// Common fields
	Timestamp   string `json:"timestamp,omitempty"`

	// For RSA types (Transactions RSA, Token)
	PrivateKey  string `json:"privateKey,omitempty"`

	// For HMAC type only
	AccessToken string `json:"accessToken,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
}

func (r *SignatureRequest) Validate() error {
	// SignatureRequestType is always required
	if r.SignatureRequestType == "" {
		return &ValidationError{Field: "signatureRequestType", Message: "signatureRequestType is required"}
	}

	switch r.SignatureRequestType {
	case SignatureTypeTransactionsRSA:
		// Requires: Method, URL, Body, PrivateKey
		if r.Method == "" {
			return &ValidationError{Field: "method", Message: "method is required for transactions"}
		}
		if r.URL == "" {
			return &ValidationError{Field: "url", Message: "url is required for transactions"}
		}
		if r.Body == nil {
			return &ValidationError{Field: "body", Message: "body is required for transactions"}
		}
		if r.PrivateKey == "" {
			return &ValidationError{Field: "privateKey", Message: "privateKey is required for RSA signatures"}
		}

	case SignatureTypeTransactionsHMAC:
		// Requires: Method, URL, Body, AccessToken, SecretKey
		if r.Method == "" {
			return &ValidationError{Field: "method", Message: "method is required for transactions"}
		}
		if r.URL == "" {
			return &ValidationError{Field: "url", Message: "url is required for transactions"}
		}
		if r.Body == nil {
			return &ValidationError{Field: "body", Message: "body is required for transactions"}
		}
		if r.AccessToken == "" {
			return &ValidationError{Field: "accessToken", Message: "accessToken is required for HMAC signatures"}
		}
		if r.SecretKey == "" {
			return &ValidationError{Field: "secretKey", Message: "secretKey is required for HMAC signatures"}
		}

	case SignatureTypeToken:
		// Requires: ClientID, PrivateKey
		if r.ClientID == "" {
			return &ValidationError{Field: "clientID", Message: "clientID is required for token generation"}
		}
		if r.PrivateKey == "" {
			return &ValidationError{Field: "privateKey", Message: "privateKey is required for token generation"}
		}

	default:
		return &ValidationError{Field: "signatureRequestType", Message: "invalid signature type"}
	}

	return nil
}

// SignatureResponse represents the response with signature components
type SignatureResponse struct {
	Signature    string            `json:"signature"`
	Timestamp    string            `json:"timestamp"`
	StringToSign string            `json:"stringToSign"`
	Headers      map[string]string `json:"headers"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type SignatureType string

const (
	SignatureTypeTransactionsRSA  SignatureType = "TRANSACTIONS_RSA_SHA256"
	SignatureTypeTransactionsHMAC SignatureType = "TRANSACTIONS_HMAC_SHA512"
	SignatureTypeToken            SignatureType = "TOKEN_RSA_SHA256"
)
