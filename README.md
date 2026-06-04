# Daytrack

FOSS daily activity tracker. Track anything you do — coding hours, workouts, reading, habits — via API or web UI. Built with Go/Gin/SQLite backend and Vite+React frontend.

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 22+
- npm

### 1. Generate JWT keys

```bash
openssl genpkey -algorithm ed25519 -out private.pem
openssl pkey -in private.pem -pubout -out public.pem
# base64 encode the keys for env vars (or inline with \n)
```

### 2. Configure environment

```bash
cp service/.env.example service/.env
```

Edit `service/.env` — at minimum set:
- `OPAQUE_SALT_NONCE` — `openssl rand -hex 16` output

JWT keys are optional — ephemeral ones are generated at startup (tokens will break on restart). For persistent tokens across restarts, set:
- `JWT_PRIVATE_KEY` — Ed25519 private key PEM (single line, literal `\n`)
- `JWT_PUBLIC_KEY` — Ed25519 public key PEM (same format)

### 3. Build frontend

```bash
cd client
npm install && npm run build
cd ..
```

### 4. Run backend

```bash
cd service
go run .
```

Open http://localhost:3000

## Docker

```bash
docker compose up --build
```

Set env vars for JWT keys via `.env` or shell environment:

```bash
JWT_PRIVATE_KEY='...' JWT_PUBLIC_KEY='...' docker compose up --build
```

## Project Structure

```
daytrack/
├── client/                    # Vite + React frontend
│   ├── src/
│   │   ├── api/
│   │   │   ├── client.js      # API client (fetch wrapper)
│   │   │   └── auth.js        # OPAQUE auth (Ed25519 via @noble/ed25519)
│   │   ├── components/
│   │   │   ├── Heatmap.jsx    # GitHub-style contribution heatmap
│   │   │   ├── Heatmap.css
│   │   │   └── Layout.jsx     # Sidebar nav + protected routes
│   │   ├── pages/
│   │   │   ├── LoginPage.jsx
│   │   │   ├── SignupPage.jsx
│   │   │   ├── DashboardPage.jsx  # Track CRUD, event tracking, heatmap
│   │   │   └── ApiKeysPage.jsx    # API key management + curl examples
│   │   └── App.jsx
│   └── package.json
├── service/                   # Go/Gin backend
│   ├── main.go
│   ├── src/
│   │   ├── api/
│   │   │   ├── server.go      # Gin server, CORS, static serving, SPA fallback
│   │   │   ├── api_key/       # API key CRUD
│   │   │   ├── event/         # Event tracking & listing
│   │   │   ├── track/         # Track CRUD
│   │   │   └── user/          # Signup/signin (OPAQUE)
│   │   ├── database/
│   │   │   ├── db.go          # SQLite init (WAL mode)
│   │   │   ├── event.go
│   │   │   ├── track.go
│   │   │   ├── api_key.go
│   │   │   └── user.go
│   │   ├── models/            # Data models, config, error types
│   │   ├── utils/             # Logging, JWT, random, middleware
│   │   └── middleware/        # Auth middleware (JWT + API key)
│   ├── .env.example
│   └── go.mod
├── Dockerfile                 # Multi-stage build (Node → Go → Debian)
├── docker-compose.yml
└── README.md
```

## Config Reference

| Variable | Default | Description |
|---|---|---|
| `HTTP_HOST` | `localhost` | Bind address |
| `HTTP_PORT` | `3000` | Listen port |
| `DB_PATH` | `./archive/database.db` | SQLite file path |
| `ENVIRONMENT` | `development` | `development` or `production` |
| `OPAQUE_SALT_NONCE` | `change-me` | Salt for OPAKE key derivation |
| `OPAQUE_CHALLENGE_RANGE` | `100000` | Challenge space (higher = more entropy) |
| `JWT_PRIVATE_KEY` | (ephemeral) | Ed25519 private key PEM |
| `JWT_PUBLIC_KEY` | (ephemeral) | Ed25519 public key PEM |
| `JWT_ACCESS_TOKEN_DURATION` | `15m` | Token expiry (Go duration format) |
| `RATE_LIMIT_BURST` | `200` | Max requests per interval per IP for event creation |
| `RATE_LIMIT_INTERVAL` | `1s` | Rate limit window (Go duration format) |
| `DISABLE_SIGNUP` | `false` | Set to `true` to disable new user registration |

## API Usage

### Authentication

Uses OPAQUE (Oblivious Password Authentication) with Ed25519 signatures. The frontend handles the challenge-response flow automatically. To use the API programmatically, you need an API key.

### API Keys

Create an API key from the web UI (API Keys page). The key is shown only once — copy it immediately.

The key is passed as a query parameter on all event endpoints:

```bash
# Track an event
curl -X POST "http://localhost:3000/v1/events/{username}/{track_name}?key={api_key}&quantity=1"

# Track with custom timestamp (RFC 3339)
curl -X POST "http://localhost:3000/v1/events/{username}/{track_name}?key={api_key}&quantity=1&created_at=2026-06-03T10:00:00Z"

# List events grouped by day
curl "http://localhost:3000/v1/events/{username}/{track_name}?key={api_key}&list_by=day"

# List raw events (each event individually)
curl "http://localhost:3000/v1/events/{username}/{track_name}?key={api_key}&list_by=raw"

# List events grouped by month
curl "http://localhost:3000/v1/events/{username}/{track_name}?key={api_key}&list_by=month"
```

### Track Management (requires JWT token)

```bash
# List tracks
curl -H "Authorization: Bearer {jwt}" http://localhost:3000/v1/tracks

# Create track
curl -X POST -H "Authorization: Bearer {jwt}" -H "Content-Type: application/json" \
  -d '{"name":"my-track"}' http://localhost:3000/v1/tracks

# Delete track
curl -X DELETE -H "Authorization: Bearer {jwt}" \
  http://localhost:3000/v1/tracks/{track_name}
```

## Tests

End-to-end tests require a running server. Start the server first, then run tests:

```bash
# Terminal 1: start server
cd service && go run .

# Terminal 2: run tests
cd service && go test ./tests/...
```

Tests cover:
- **Auth**: challenge, signup, signin, duplicate username
- **API keys**: create, list, delete
- **Tracks**: create, list, update, delete
- **Events**: create at time, list by day/public/private
- **Access guards**: anonymous vs authenticated access for private, public-read, and public-read-write tracks
- **Cross-access**: users accessing other users' tracks with different visibility levels

## Development

### Frontend dev server (with backend proxy)

```bash
# Start backend
cd service && go run . &

# Start Vite dev server (proxy API calls to backend)
cd client && VITE_API_URL=http://localhost:3000 npm run dev
```

### Production build (single binary serves everything)

```bash
cd client && npm install && npm run build
cd ../service && go build -o daytrack .
./daytrack
# Frontend at http://localhost:3000 (served by Gin)
```

## License

MIT. See [LICENSE](LICENSE).
