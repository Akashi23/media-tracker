# Media Tracker Backend

A Go backend for the Media Tracker application that allows users to track their media consumption (movies, books, anime, games, TV shows, videos).

## Features

- **Authentication**: JWT-based authentication with email login
- **Media Management**: Create and search media items
- **Entry Tracking**: Track your progress, ratings, and reviews
- **Collections**: Create and share collections
- **Guest Mode**: Use without registration, data stored locally
- **Sharing**: Share collections and profiles via public links
- **Data Migration**: Merge guest data into registered accounts

## Tech Stack

- **Language**: Go 1.24.6+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with JSONB support
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Logging**: Zerolog

## Quick Start

### Prerequisites

- Go 1.24.6+
- PostgreSQL 12+
- Redis 6+

### Installation

1. Clone the repository
2. Copy environment file:
   ```bash
   cp env.example .env
   ```

3. Update `.env` with your database and Redis credentials

4. Run database migration:
   ```bash
   psql -U your_user -d media_tracker -f ../migrations/001_initial_schema.sql
   ```

5. Install dependencies:
   ```bash
   go mod tidy
   ```

6. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login with email
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get user profile

### Media
- `POST /api/media` - Create media item
- `GET /api/media/search?q=query&type=movie` - Search media

### Entries
- `GET /api/entries` - List user entries
- `POST /api/entries` - Create entry
- `GET /api/entries/:id` - Get entry
- `PATCH /api/entries/:id` - Update entry
- `DELETE /api/entries/:id` - Delete entry

### Collections
- `GET /api/collections` - List collections
- `POST /api/collections` - Create collection
- `GET /api/collections/:id` - Get collection
- `PATCH /api/collections/:id` - Update collection
- `DELETE /api/collections/:id` - Delete collection
- `POST /api/collections/:id/share` - Create share link

### Guest
- `POST /api/guest/snapshot` - Create guest data snapshot
- `POST /api/guest/merge` - Merge guest data to account

### Public
- `GET /s/:token` - View public share

## Docker

Build and run with Docker:

```bash
docker build -t media-tracker-backend .
docker run -p 8080:8080 --env-file .env media-tracker-backend
```

## Development

### Project Structure

```
back/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connections
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Middleware functions
│   ├── models/            # Data models and DTOs
│   ├── repository/        # Data access layer
│   └── services/          # Business logic
├── migrations/            # Database migrations
└── Dockerfile            # Docker configuration
```

### Adding New Features

1. Define models in `internal/models/`
2. Create repository methods in `internal/repository/`
3. Implement business logic in `internal/services/`
4. Add HTTP handlers in `internal/handlers/`
5. Register routes in `main.go`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `media_tracker` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | - |
| `REDIS_DB` | Redis database | `0` |
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_EXPIRY` | JWT expiry hours | `24` |

## License

MIT
