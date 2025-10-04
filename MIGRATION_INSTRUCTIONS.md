# Migration Instructions for Hourly Worth Feature

To complete the implementation of the hourly worth feature, you need to run the following commands:

## 1. Run Database Migration

First, ensure your database is running and accessible, then run:

```bash
make migrate-up
```

Or if make is not available:
```bash
dbmate up
```

This will apply the migration file `db/migrations/20251006_add_hourly_worth_to_users.sql` which adds the `hourly_worth` column to the users table.

## 2. Generate SQL Code

After the migration is complete, generate the Go code from the SQL queries:

```bash
make sqlc-generate
```

Or if make is not available:
```bash
sqlc generate
```

This will update the `internal/database/users.sql.go` file with the new methods:
- `UpdateUserHourlyWorth`
- `GetUserHourlyWorth`

## 3. Test the Implementation

1. Start the backend server:
   ```bash
   make run
   # or
   go run main.go
   ```

2. Start the frontend development server:
   ```bash
   cd web
   npm run dev
   ```

3. Navigate to the settings page and test:
   - The hourly worth should be fetched from the database on page load
   - Changing the value and clicking "Save Changes" should update it in the database
   - The value should persist across page refreshes

## Implementation Summary

The following changes have been made:

### Backend:
- Created migration to add `hourly_worth` column to users table
- Added SQL queries for getting and updating hourly worth
- Created `SettingsHandler` with endpoints:
  - `GET /api/settings/hourly-worth` - Get user's hourly worth
  - `PUT /api/settings/hourly-worth` - Update user's hourly worth
- Added routes to main.go

### Frontend:
- Added API methods to `api.ts` for hourly worth operations
- Created `+page.ts` to load hourly worth from API on page load
- Updated settings page to:
  - Display fetched hourly worth value
  - Save changes to the database via API
  - Show loading and error states
  - Sync with local store for consistency

The hourly worth is now stored in the database and synchronized between the frontend and backend.