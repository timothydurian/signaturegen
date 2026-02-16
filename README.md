# SNAP Signature Generator

A Go-based HTTP service for generating cryptographic signatures for SNAP (Standard National Application Programming Interface) API authentication. Supports multiple signature types including RSA-SHA256 and HMAC-SHA512.

## Features

- 🔐 **Multiple Signature Types**
  - RSA-SHA256 for transactions
  - RSA-SHA256 for token generation
  - HMAC-SHA512 for transactions

- 🚀 **Flexible Body Handling**
  - Accepts JSON objects
  - Accepts raw string bodies (ideal for Postman/API testing)

- ✅ **Production Ready**
  - Input validation
  - Comprehensive error handling
  - RESTful API design

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Signature Formats](#signature-formats)
- [Examples](#examples)
- [Postman Integration](#postman-integration)

## Installation

### Prerequisites

- Go 1.21 or higher (for local development)
- Docker & Docker Compose (for containerized deployment)

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd signaturegen

# Using docker-compose (easiest)
docker-compose up -d

# Or using Makefile
make docker-up

# Check logs
docker-compose logs -f
# or
make docker-logs
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd signaturegen

# Install dependencies
go mod download

# Build the application
go build
# or
make build

# Run the service
./signaturegen
# or
make run
```

The service will start on `http://localhost:8080` by default.

### Using Makefile

The project includes a Makefile for common operations:

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run locally
make test              # Run tests
make docker-build      # Build Docker image
make docker-up         # Start with docker-compose
make docker-down       # Stop docker services
make docker-logs       # View logs
make docker-rebuild    # Rebuild and restart
```

## Usage

### Quick Start

```bash
# Start the service
./signaturegen

# Generate a signature (example using curl)
curl -X POST http://localhost:8080/generate-signature \
  -H "Content-Type: application/json" \
  -d '{
    "signatureRequestType": "TRANSACTIONS_HMAC_SHA512",
    "method": "POST",
    "url": "/api/v1/transfer",
    "body": "{\"amount\":\"10000\"}",
    "accessToken": "your-b2b-access-token",
    "secretKey": "your-secret-key"
  }'
```

## API Reference

### Generate Signature

**Endpoint:** `POST /generate-signature`

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `signatureRequestType` | string | Yes | Type of signature: `TRANSACTIONS_RSA_SHA256`, `TRANSACTIONS_HMAC_SHA512`, or `TOKEN_RSA_SHA256` |
| `method` | string | Conditional | HTTP method (required for transaction signatures) |
| `url` | string | Conditional | Relative URL path (required for transaction signatures) |
| `body` | string/object | Conditional | Request body - can be JSON object or raw string (required for transaction signatures) |
| `clientID` | string | Conditional | Client ID (required for token signatures) |
| `timestamp` | string | No | ISO 8601 timestamp (auto-generated if not provided) |
| `privateKey` | string | Conditional | RSA private key in PEM format (required for RSA signatures) |
| `accessToken` | string | Conditional | B2B access token (required for HMAC signatures) |
| `secretKey` | string | Conditional | Secret key (required for HMAC signatures) |

**Response:**

```json
{
  "signature": "base64-encoded-signature",
  "timestamp": "2024-02-16T10:30:00.000+07:00",
  "stringToSign": "POST:/api/v1/transfer:access-token:hash:timestamp",
  "headers": {
    "X-TIMESTAMP": "2024-02-16T10:30:00.000+07:00",
    "X-SIGNATURE": "base64-encoded-signature"
  }
}
```

### Health Check

**Endpoint:** `GET /health`

Returns service health status.

## Signature Formats

### 1. Transaction HMAC-SHA512

**Format:**
```
<HTTP_METHOD>:<RELATIVE_PATH_URL>:<B2B_ACCESS_TOKEN>:<BODY_HASH>:<TIMESTAMP>
```

**Body Hash Calculation:**
```
LowerCase(HexEncode(SHA-256(Minify(<HTTP_BODY>))))
```

**Example String to Sign:**
```
POST:/api/v1/transfer:eyJhbGc...:a1b2c3d4e5f6:2024-02-16T10:30:00.000+07:00
```

### 2. Transaction RSA-SHA256

**Format:**
```
<HTTP_METHOD>:<RELATIVE_PATH_URL>:<BODY_HASH>:<TIMESTAMP>
```

**Example String to Sign:**
```
POST:/api/v1/transfer:a1b2c3d4e5f6:2024-02-16T10:30:00.000+07:00
```

### 3. Token RSA-SHA256

**Format:**
```
<CLIENT_ID>|<TIMESTAMP>
```

**Example String to Sign:**
```
your-client-id|2024-02-16T10:30:00.000+07:00
```

## Examples

### Example 1: HMAC Transaction Signature

```json
{
  "signatureRequestType": "TRANSACTIONS_HMAC_SHA512",
  "method": "POST",
  "url": "/api/v1/debit/payment-host-to-host",
  "body": "{\"partnerReferenceNo\":\"12345\",\"amount\":{\"value\":\"10000.00\",\"currency\":\"IDR\"}}",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "secretKey": "your-secret-key-here"
}
```

**Response:**
```json
{
  "signature": "xVn8kPqR7mN4...",
  "timestamp": "2024-02-16T10:30:00.000+07:00",
  "stringToSign": "POST:/api/v1/debit/payment-host-to-host:eyJhbG...:7f3d8a2b1c:2024-02-16T10:30:00.000+07:00",
  "headers": {
    "X-TIMESTAMP": "2024-02-16T10:30:00.000+07:00",
    "X-SIGNATURE": "xVn8kPqR7mN4..."
  }
}
```

### Example 2: RSA Transaction Signature

```json
{
  "signatureRequestType": "TRANSACTIONS_RSA_SHA256",
  "method": "POST",
  "url": "/api/v1/transfer",
  "body": {
    "amount": "10000",
    "currency": "IDR"
  },
  "privateKey": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASC...\n-----END PRIVATE KEY-----"
}
```

### Example 3: Token Signature

```json
{
  "signatureRequestType": "TOKEN_RSA_SHA256",
  "clientID": "your-client-id",
  "privateKey": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASC...\n-----END PRIVATE KEY-----"
}
```

## Postman Integration

### Pre-Request Script

Use this script in Postman to automatically generate signatures before each request:

```javascript
console.log("=== SNAP HMAC SIGNATURE PRE-REQUEST START ===");

// Configuration
const SIGN_URL = pm.environment.get("SIGNATURE_SERVICE_URL");
const CLIENT_SECRET = pm.environment.get("SNAP_SECRET_KEY");
const b2bAccessToken = pm.environment.get("SNAP_B2B_TOKEN");

// Request data
const method = pm.request.method;
const url = pm.request.url.getPath();
const rawBody = pm.request.body ? pm.request.body.raw : null;

// Build payload
const payload = {
    signatureRequestType: "TRANSACTIONS_HMAC_SHA512",
    method: method,
    url: url,
    body: rawBody,  // Raw string body
    accessToken: b2bAccessToken,
    secretKey: CLIENT_SECRET
};

console.log("Signature Payload:", payload);

// Call signature service
pm.sendRequest({
    url: SIGN_URL,
    method: "POST",
    header: {
        "Content-Type": "application/json"
    },
    body: {
        mode: "raw",
        raw: JSON.stringify(payload)
    }
}, function (err, res) {
    if (err) {
        console.log("Signature request error:", err);
        return;
    }

    const json = res.json();
    console.log("STRING TO SIGN:", json.stringToSign);
    console.log("TIMESTAMP:", json.timestamp);
    console.log("SIGNATURE:", json.signature);

    // Set headers for the main request
    pm.environment.set("X-SIGNATURE", json.headers["X-SIGNATURE"]);
    pm.environment.set("X-TIMESTAMP", json.headers["X-TIMESTAMP"]);

    console.log("=== SNAP HMAC SIGNATURE PRE-REQUEST END ===");
});
```

### Required Environment Variables

Set these variables in your Postman environment:

| Variable | Description |
|----------|-------------|
| `SIGNATURE_SERVICE_URL` | URL of this signature service (e.g., `http://localhost:8080/generate-signature`) |
| `SNAP_SECRET_KEY` | Your SNAP secret key |
| `SNAP_B2B_TOKEN` | Your B2B access token |

## Error Handling

### Validation Errors

```json
{
  "error": "method is required for transactions"
}
```

### Common Error Scenarios

- Missing required fields
- Invalid signature type
- Invalid private key format
- Malformed JSON body

## Security Considerations

⚠️ **Important Security Notes:**

1. **Never commit private keys** to version control
2. **Use environment variables** for sensitive data
3. **Rotate keys regularly**
4. **Use HTTPS** in production
5. **Implement rate limiting** for production deployments
6. **Store secret keys securely** (use secrets management tools)

## Docker Deployment

### Quick Start with Docker

```bash
# Start the service
docker-compose up -d

# Check if it's running
curl http://localhost:8080/health

# View logs
docker-compose logs -f signaturegen

# Stop the service
docker-compose down
```

### Docker Configuration

The Docker setup includes:
- **Multi-stage build** for optimized image size
- **Health checks** for container monitoring
- **Resource limits** (CPU: 0.5, Memory: 256MB)
- **Auto-restart** policy
- **Alpine-based** minimal image

### Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
cp .env.example .env
```

Available environment variables:
- `PORT` - Server port (default: 8080)
- `GO_ENV` - Environment (development/production)
- `LOG_LEVEL` - Logging level (info/debug/error)

### Production Deployment

For production, consider:

1. **Use environment variables** for secrets
2. **Enable HTTPS** with a reverse proxy (nginx/traefik)
3. **Set up monitoring** (Prometheus/Grafana)
4. **Configure logging** to external service
5. **Adjust resource limits** based on load

Example with nginx reverse proxy:

```nginx
server {
    listen 443 ssl;
    server_name signature.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Using make
make test
```

### Development with Hot Reload

Install [Air](https://github.com/cosmtrek/air) for hot reload:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Project Structure

```
signaturegen/
├── main.go              # HTTP server and routes
├── handlers.go          # HTTP handlers
├── signer.go            # Signature generation logic
├── utils.go             # Utility functions (SHA-256, RSA, HMAC)
├── models.go            # Request/response models
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose configuration
├── .dockerignore        # Docker ignore patterns
├── .env.example         # Environment variables template
├── Makefile             # Build and deployment commands
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[Add your license here]

## Support

For issues and questions:
- Create an issue in the repository
- Contact: [your-email@example.com]

## Acknowledgments

- Implements SNAP (Standard National Application Programming Interface) specification
- Compatible with Indonesian payment gateway standards
