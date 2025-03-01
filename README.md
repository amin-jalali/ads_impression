# Ad Impression Counter Service

This repository contains an **in-memory concurrent service** for tracking ad impressions in real-time. The service:

- **Registers advertising campaigns** (with name and start time)
- **Tracks impressions** for each campaign
- **Handles duplicate impressions** via a TTL mechanism
- **Provides stats** (LastHour, LastDay, TotalCount) for each campaign
- **Exposes a REST API** on standard HTTP endpoints
- **Runs concurrency-safe** using Go's `sync.Mutex`.

---

## Features

1. **Campaign Management**
    - Create a new campaign with `name` and `start_time`.
    - Persists the campaign in memory using a `Server` struct.
2. **Impression Tracking**
    - Tracks user impressions by `campaign_id`, `user_id`, and `ad_id`.
    - Uses a TTL (1 hour) to avoid counting **duplicate** impressions from the same user.
3. **Campaign Stats**
    - Retrieves aggregated stats for a given campaign.
    - Stats include impressions in the **Last Hour**, **Last Day**, and **Total**.
4. **API Endpoints**
    - `POST /api/v1/campaigns` — Create a campaign
        - Request Body: `{"name": "string", "start_time": "RFC3339 Time"}`
    - `POST /api/v1/impressions` — Track an impression
        - Request Body: `{"campaign_id": "string", "user_id": "string", "ad_id": "string"}`
    - `GET /api/v1/campaigns/{id}/stats` — Get stats by campaign ID
    - `404` handler for invalid routes
5. **Concurrent & Thread-Safe**
    - Uses `sync.Mutex` to lock shared resources.
6. **Test Coverage**
    - Includes an extensive suite of unit tests under `internal/handlers/tests`.

---

## Project Structure

```
.
├── cmd
│   └── server
│       ├── config.go           # Configuration logic (if any)
│       └── main.go             # Entry point for running the service
├── docker-compose.yml          # Docker Compose setup
├── Dockerfile                  # Docker build file
├── go.mod
├── go.sum
├── internal
│   ├── handlers
│   │   ├── campaign.go         # Create campaign handler
│   │   ├── impression.go       # Track impression handler
│   │   ├── notFound.go         # NotFoundHandler for invalid routes
│   │   ├── server.go           # Defines Server struct, shared concurrency-safe maps
│   │   ├── stats.go            # Get campaign stats handler
│   │   └── tests
│   │       ├── campaign_test.go
│   │       ├── config_test.go
│   │       ├── high_volume_test.go
│   │       ├── impression_test.go
│   │       ├── main_test.go
│   │       ├── not_found_test.go
│   │       ├── server_test.go
│   │       └── stats_test.go
│   └── logger
│       └── logger.go           # Custom logger logic
└── main.go                     # (Optional) top-level main, or empty if using cmd/server
```

### Key Components

- ``: Minimal entry point that sets up and starts the HTTP server.
- ``: Business logic for creating campaigns, tracking impressions, retrieving stats, and handling 404s.
- ``: Defines the `Server` struct, storing concurrency-safe maps for campaigns, impressions, and stats.
- ``: Logging utilities.
- ``** & **``: For containerization if desired.
- ``: Houses the test files for campaigns, impressions, stats, concurrency, etc.

---

## Getting Started

### Prerequisites

- [Go 1.18+](https://golang.org/dl/)
- (Optional) [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone** this repository:

```bash
git clone https://github.com/yourusername/ad-impression-service.git
cd ad-impression-service
```

2. **Install dependencies**:

```bash
go mod download
```

### Running Locally

#### Via Go:

```bash
go run ./cmd/server
```

The server starts on `:8080`. You should see:

```
Server started on :8080
```

#### Via Docker Compose:

1. Make sure Docker is running.
2. Run:
   ```bash
   docker-compose up --build
   ```
3. The service will be available on `:8080` (or whichever port configured in `docker-compose.yml`).

---

## Usage

Below are sample requests for each endpoint. You can use [curl](https://curl.se/) or a tool like [Postman](https://www.postman.com/).

### 1. Create a Campaign

```bash
curl -X POST \
  http://localhost:8080/api/v1/campaigns \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Campaign A",
    "start_time": "2025-01-01T00:00:00Z"
}'
```

#### Example Response

```json
{
  "id": "some-uuid-value",
  "name": "Campaign A",
  "start_time": "2025-01-01T00:00:00Z"
}
```

### 2. Track an Impression

```bash
curl -X POST \
  http://localhost:8080/api/v1/impressions \
  -H 'Content-Type: application/json' \
  -d '{
    "campaign_id": "some-uuid-value",
    "user_id": "user123",
    "ad_id": "ad456"
}'
```

#### Example Response

```json
{ "message": "Impression saved successfully" }
```

### 3. Get Campaign Stats

```bash
curl -X GET http://localhost:8080/api/v1/campaigns/some-uuid-value
```

#### Example Response

```json
{
  "campaign_id": "some-uuid-value",
  "last_hour": 10,
  "last_day": 50,
  "total": 100
}
```

### 4. Unknown Routes

Any route not registered returns `404 not found`.

---

## Testing

```bash
go test -coverprofile=coverage.out ./...
```

Check coverage:

```bash
go tool cover -func=coverage.out
```

Generate an HTML report:

```bash
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

---

## Contributing

1. **Fork** the repository.
2. **Create** a feature branch: `git checkout -b feature/new-feature`
3. **Commit** your changes: `git commit -m 'Add some feature'`
4. **Push** to the branch: `git push origin feature/new-feature`
5. **Open** a Pull Request.


