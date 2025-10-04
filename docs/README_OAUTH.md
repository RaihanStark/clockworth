# Google OAuth Setup Guide

## Prerequisites

1. Install required tools:
   - Docker & Docker Compose
   - Go 1.24+
   - dbmate: `brew install dbmate` or download from https://github.com/amacneil/dbmate
   - sqlc: `brew install sqlc` or download from https://sqlc.dev

## Setup Steps

### 1. Start PostgreSQL

```bash
make dev-up
# or
docker-compose up -d
```

### 2. Configure Environment Variables

Edit `.env` file with your Google OAuth credentials:

```bash
GOOGLE_CLIENT_ID=your_google_client_id_here
GOOGLE_CLIENT_SECRET=your_google_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
SESSION_SECRET=your_session_secret_here_min_32_chars
DATABASE_URL=postgres://clockworth:clockworth_dev_password@localhost:5432/clockworth?sslmode=disable
FRONTEND_URL=http://localhost:5173
PORT=8080
```

### 3. Get Google OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Go to Credentials → Create Credentials → OAuth 2.0 Client ID
5. Configure OAuth consent screen
6. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
7. Copy Client ID and Client Secret to `.env`

### 4. Run Migrations

```bash
make migrate-up
# or
dbmate up
```

### 5. Generate SQLC Code

```bash
make sqlc-generate
# or
sqlc generate
```

### 6. Run the Application

```bash
make run
# or
go run main.go
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /auth/google/login` - Start OAuth flow, returns login URL
- `GET /auth/google/callback` - OAuth callback (redirects to frontend)

### Protected Endpoints (requires authentication)

- `GET /api/me` - Get current user info
- `POST /api/logout` - Logout current user

## Authentication Flow

1. Frontend calls `GET /auth/google/login` to get OAuth URL
2. Redirect user to the OAuth URL
3. User authenticates with Google
4. Google redirects to `GET /auth/google/callback`
5. Backend creates session and redirects to frontend with session token
6. Frontend stores token and includes it in subsequent requests

### Using the Token

Include token in requests either:

**Header:**
```
Authorization: Bearer <token>
```

**Cookie:**
```
session_token=<token>
```

## Database Schema

### Users Table
- `id` - UUID primary key
- `email` - User email (unique)
- `name` - User full name
- `picture` - Profile picture URL
- `google_id` - Google user ID (unique)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Sessions Table
- `id` - UUID primary key
- `user_id` - Foreign key to users
- `token` - Session token (unique)
- `expires_at` - Expiration timestamp
- `created_at` - Creation timestamp

## Development Commands

```bash
# Start database
make dev-up

# Stop database
make dev-down

# Run migrations
make migrate-up

# Rollback migration
make migrate-down

# Create new migration
make migrate-create name=add_new_table

# Generate SQLC code
make sqlc-generate

# Run application
make run

# Build application
make build
```

## Testing the OAuth Flow

### Backend Only (API Testing)
1. Start the backend: `make run`
2. Visit: http://localhost:8080/auth/google/login
3. You'll receive JSON with a `url` field
4. Open that URL in browser to authenticate with Google
5. After authentication, you'll be redirected to your frontend with the session token

### Full Stack (Frontend + Backend)

1. **Start Backend:**
   ```bash
   make run
   # or
   go run main.go
   ```

2. **Start Frontend:**
   ```bash
   cd web
   npm install  # if not done yet
   npm run dev
   ```

3. **Test the Flow:**
   - Open http://localhost:5173/login in your browser
   - Click "Sign in with Google"
   - Authenticate with Google
   - You'll be redirected to `/callback`, then to `/dashboard`
   - Your session is now active!

4. **Verify Authentication:**
   - Navigate to http://localhost:5173/dashboard
   - You should see your dashboard (protected route)
   - If not logged in, you'll be redirected to login

### Frontend Environment Variables

The frontend uses `.env` file in the `web/` directory:

```bash
VITE_API_URL=http://localhost:8080
```

This tells the frontend where the backend API is running.

## Troubleshooting

- **Database connection error**: Ensure PostgreSQL is running (`make dev-up`)
- **Migration errors**: Check `.dbmate` file has correct DATABASE_URL
- **OAuth errors**: Verify Google credentials in `.env` are correct
- **CORS errors**: Update `FRONTEND_URL` in `.env` to match your frontend URL
