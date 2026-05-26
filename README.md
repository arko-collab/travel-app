# Travel App AI Backend

This project is a Go HTTP API for a travel workflow that covers three main steps:

- extracting travel intent from free-form text with Gemini,
- searching an in-memory catalog of trip bundles,
- creating approval requests backed by PostgreSQL and notifying an approver by email via SendGrid.

The server exposes JSON APIs under `/api/*` and uses Gorilla Mux, Viper, GORM, PostgreSQL, and SendGrid.

## What the Project Does

The backend is organized around a simple travel request flow:

1. A user sends natural language text to `/api/intent`.
2. The service extracts structured travel details such as destination, dates, and purpose.
3. The client sends a destination to `/api/search` to retrieve available trip bundles.
4. The selected trip can be submitted to `/api/approval`.
5. The approval request is stored in PostgreSQL and an approval email is sent to the configured approver.

## Features

- Intent extraction using Gemini with internal fallback behavior if the API call fails.
- Search endpoint backed by seeded travel bundles for quick local testing.
- Approval workflow persisted in PostgreSQL.
- SendGrid email notification for new approval requests.
- Graceful HTTP server shutdown.
- Request logging, JSON responses, and permissive CORS middleware.

## Tech Stack

- Go
- Gorilla Mux
- Viper
- GORM
- PostgreSQL
- SendGrid
- Gemini API

## Project Structure

```text
travelapp-ai/
|-- cmd/
|   `-- server/
|       `-- main.go                  # Application entrypoint and route registration
|-- internal/
|   |-- approval/
|   |   |-- handler.go              # HTTP handler for approval submissions
|   |   |-- models.go               # Approval request/response payloads
|   |   `-- service.go              # Approval creation and SendGrid email logic
|   |-- config/
|   |   `-- config.go               # Loads .env values with Viper
|   |-- db/
|   |   `-- postgresdb.go           # PostgreSQL connection and auto-migration
|   |-- intent/
|   |   |-- handler.go              # HTTP handler for intent extraction
|   |   |-- models.go               # Intent request/response payloads
|   |   `-- service.go              # Gemini integration and fallback parsing
|   |-- middleware/
|   |   `-- logging.go              # Logging, JSON response helpers, CORS middleware
|   `-- search/
|       |-- handler.go              # HTTP handler for bundle search
|       |-- models.go               # Search request and bundle response models
|       `-- service.go              # In-memory seeded search catalog
|-- tmp/                            # Runtime/temp artifacts
|-- .air.toml                       # Air live-reload configuration for development
|-- .env                            # Local environment file loaded at startup
|-- .env.example                    # Safe template for local setup
|-- go.mod                          # Go module definition
|-- go.sum                          # Dependency checksums
`-- Travel-app.postman_collection.json
                                 # Postman collection for local API testing
```

## Environment Setup

The application reads configuration from a `.env` file in the repository root.

### 1. Create the env file

Copy `.env.example` to `.env` and replace the placeholder values with real credentials.

```powershell
Copy-Item .env.example .env
```

### 2. Required variables

| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | Yes | PostgreSQL connection string used by GORM |
| `SENDGRID_API_KEY` | Yes | SendGrid API key for approval emails |
| `SENDGRID_FROM` | Yes | Sender email address used in approval emails |
| `APPROVER_EMAIL` | Yes | Recipient email for approval notifications |
| `GEMINI_API_KEY` | Yes | API key for Gemini intent extraction |
| `PORT` | No | HTTP server port, defaults to `8080` |

### 3. Example `.env`

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/travelapp?sslmode=disable
SENDGRID_API_KEY=your_sendgrid_api_key
SENDGRID_FROM=no-reply@example.com
APPROVER_EMAIL=manager@example.com
GEMINI_API_KEY=your_gemini_api_key
PORT=8080
```

### 4. Environment file placement

Keep the file at the project root:

```text
travelapp-ai/
|-- .env
|-- .env.example
|-- cmd/
|-- internal/
`-- go.mod
```

Do not place `.env` under `cmd/` or `internal/`. The current config loader expects `.env` in the repository root.

## Prerequisites

Before running the app, make sure you have:

- a Go toolchain compatible with the version declared in `go.mod`,
- a running PostgreSQL instance,
- a SendGrid account and API key,
- a Gemini API key.

## Local Development

### Install dependencies

Go will fetch dependencies automatically when you build or run the project.

```powershell
go mod download
```

### Run the server

From the repository root:

```powershell
go run ./cmd/server
```

The API starts on `http://localhost:8080` unless `PORT` is overridden.

### Optional: live reload with Air

The repository contains `.air.toml` for live reload. If you use Air:

```powershell
air
```

## Database Notes

- The server connects to PostgreSQL on startup.
- The `approval_requests` table is auto-migrated by GORM.
- If the database connection fails, the server exits during startup.

## API Endpoints

### Health check

`GET /api/health`

Response:

```json
{
  "status": "ok"
}
```

### Extract intent

`POST /api/intent`

Request:

```json
{
  "text": "I want to book a flight from Lahore to Singapore on next monday"
}
```

Response shape:

```json
{
  "destination": "Singapore, Singapore",
  "dateFrom": "2026-05-25",
  "dateTo": "2026-05-26",
  "purpose": "Business Travel"
}
```

### Search trips

`POST /api/search`

Request:

```json
{
  "destination": "Singapore, Singapore"
}
```

Response shape:

```json
[
  {
    "id": 1,
    "label": "Best value",
    "flight": {
      "airline": "Lufthansa",
      "code": "LH203",
      "dep": "Kolkata",
      "arr": "Berlin",
      "type": "Economy"
    },
    "hotel": {
      "name": "Hilton Berlin",
      "stars": 4,
      "breakfast": true
    },
    "price": 487,
    "co2": 142,
    "durationMin": 165,
    "inPolicy": true,
    "policyNote": ""
  }
]
```

### Submit approval request

`POST /api/approval`

Request:

```json
{
  "destination": "Berlin, Germany",
  "dateFrom": "2026-05-28",
  "dateTo": "2026-05-29",
  "purpose": "Client Meeting",
  "flightInfo": "Lufthansa LH203",
  "hotelInfo": "Hilton Berlin",
  "totalCost": 487,
  "notes": "Need manager approval"
}
```

Response shape:

```json
{
  "approvalId": "APR-1748390000000000000",
  "status": "PENDING",
  "message": "Approval request submitted successfully"
}
```

## Testing the API

Use the included Postman collection:

- `Travel-app.postman_collection.json`

Set the collection variable `server` to your local base URL, for example:

```text
http://localhost:8080
```

## Suggested Developer Workflow

1. Start PostgreSQL locally.
2. Copy `.env.example` to `.env`.
3. Fill in the required API keys and email addresses.
4. Run `go run ./cmd/server`.
5. Call `/api/health`.
6. Test `/api/intent`, `/api/search`, and `/api/approval` with Postman.

## Notes

- `/api/search` currently uses seeded in-memory data rather than a live provider.
- `/api/approval` stores requests in PostgreSQL and then sends the email asynchronously.
- The config loader fails fast if any required environment variables are missing.

