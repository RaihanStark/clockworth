# Clockworth ⏳

> Personal time economics and activity ROI tracking application.

Clockworth helps you understand the monetary and productive value of your time. Log daily activities, assign ROI classifications, and measure your daily output against your target hourly rate.

---

## Features

- **Interactive Timeline**: Visual 24-hour daily time-blocking grid to schedule, view, edit, and delete activities.
- **Time Economics & Hourly Worth**: Set your personal hourly worth rate and calculate the monetary equivalent of tracked blocks.
- **ROI Categorization**: Tag activities with ROI levels (`high`, `medium`, `low`) to see where your highest leverage hours go.
- **30-Day Activity Heatmap**: Visualize activity consistency and productivity patterns over the past month.
- **Dual Authentication**:
  - Secure **Email & Password** authentication with bcrypt hashing.
  - **Google OAuth2** single sign-on with HTTP-only cookie session management.
- **Responsive Web Interface**: Modern UI built with SvelteKit and TailwindCSS.

---

## Tech Stack

- **Backend**: Go 1.24, Echo v4 web framework
- **Database**: PostgreSQL with [sqlc](https://sqlc.dev/) for type-safe SQL queries and [dbmate](https://github.com/amacneil/dbmate) for migrations
- **Frontend**: SvelteKit, TypeScript, TailwindCSS, Lucide Icons
- **Deployment**: Docker Compose for local dev, Render Blueprint (`render.yaml`) for cloud deployment

---

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+ and npm
- Docker and Docker Compose
- [dbmate](https://github.com/amacneil/dbmate) and [sqlc](https://sqlc.dev/) (optional, pre-generated queries included)

### 1. Clone & Configure Environment

```bash
git clone https://github.com/RaihanStark/clockworth.git
cd clockworth

cp .env.example .env
```

Review `.env`:
```env
# Database
DATABASE_URL=postgres://clockworth:clockworth_dev_password@localhost:5432/clockworth?sslmode=disable

# Frontend & Server
FRONTEND_URL=http://localhost:5173
PORT=8080

# Optional: Google OAuth (leave empty to use Email & Password only)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

### 2. Start PostgreSQL & Run Migrations

```bash
# Start Postgres container
make dev-up
# or: docker compose up -d

# Run migrations
make migrate-up
# or: dbmate up
```

### 3. Run the Backend

```bash
make run
# or: go run main.go
```

The API server will listen on `http://localhost:8080`.

### 4. Run the Frontend

In a separate terminal:

```bash
cd web
npm install
npm run dev
```

Open `http://localhost:5173` in your browser.

---

## API Endpoints

### Public Endpoints
- `GET /health` - Server health check
- `POST /auth/register` - Register a new account with email & password
- `POST /auth/login` - Authenticate with email & password
- `GET /auth/google/login` - Initiate Google OAuth2 flow
- `GET /auth/google/callback` - Google OAuth2 callback handler

### Protected Endpoints (Requires Session Cookie)
- `GET /api/me` - Get authenticated user profile
- `POST /api/logout` - Invalidate session and clear auth cookie
- `GET /api/settings` - Retrieve user settings (hourly worth)
- `PUT /api/settings/hourly-worth` - Update user's hourly worth
- `GET /api/activities?date=YYYY-MM-DD` - Get activities for a specific date
- `GET /api/activities/heatmap` - Get 30-day activity count heatmap
- `POST /api/activities` - Log a new activity block
- `GET /api/activities/:id` - Fetch details of an activity
- `PUT /api/activities/:id` - Update an activity
- `DELETE /api/activities/:id` - Delete an activity

---

## Development

```bash
# Start/stop database
make dev-up
make dev-down

# Migrations
make migrate-up
make migrate-down
make migrate-create name=add_something

# Code generation
make sqlc-generate

# Run tests
go test ./...
```

---

## License

MIT
