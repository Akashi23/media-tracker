# Media Tracker

A comprehensive media tracking application where users can track their consumption of movies, books, anime, games, TV shows, and videos. Features include progress tracking, ratings, reviews, collections, sharing capabilities, and guest mode for immediate use without registration.

## 🚀 Quick Start with Docker

### Prerequisites
- Docker and Docker Compose
- Git

### Production Setup
```bash
# Clone the repository
git clone <repository-url>
cd media-tracker

# Start the full stack
docker-compose up -d

# Check status
docker-compose ps
```

### Development Setup
```bash
# Start development environment with hot reload
make dev-up

# View logs
make docker-logs

# Stop development environment
make dev-down
```

## 📁 Project Structure

```
media-tracker/
├── back/                    # Go backend
│   ├── internal/           # Internal packages
│   │   ├── config/         # Configuration
│   │   ├── database/       # Database connections
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Middleware
│   │   ├── models/         # Data models
│   │   ├── repository/     # Data access layer
│   │   └── services/       # Business logic
│   ├── main.go            # Application entry point
│   ├── Dockerfile         # Production Dockerfile
│   ├── Dockerfile.dev     # Development Dockerfile
│   └── .air.toml         # Hot reload configuration
├── front/                  # SvelteKit frontend
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/  # UI components
│   │   │   ├── stores/      # State management
│   │   │   ├── types/       # TypeScript types
│   │   │   └── utils/       # Utilities
│   │   └── routes/          # SvelteKit routes
│   ├── static/              # Static assets
│   └── package.json
├── migrations/             # Database migrations
├── scripts/                # Utility scripts
├── .github/workflows/      # CI/CD workflows
├── docker-compose.yml      # Production stack
├── docker-compose.dev.yml  # Development stack
└── Makefile               # Development commands
```

## 🛠️ Development

### Local Development (without Docker)
```bash
# Install dependencies
make deps

# Run locally (requires PostgreSQL and Redis)
make run

# Run tests
make test

# Format code
make fmt
```

### Docker Development
```bash
# Start development environment
make dev-up

# View logs
make docker-logs

# Stop environment
make dev-down
```

### Database Management
```bash
# Run migrations
make db-migrate

# Add seed data (additional users and media)
make db-seed

# Reset database (drop and recreate with seed data)
make db-reset

# Connect to database
docker exec -it media-tracker-postgres-dev psql -U postgres -d media_tracker

# Add default user programmatically
cd scripts && go run add_default_user.go
```

## 🌐 Services

| Service | Port | Description |
|---------|------|-------------|
| Backend API | 8080 | Go REST API |
| Frontend | 3000 | SvelteKit web interface |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache |

## ✨ Features

### 🎯 Core Functionality
- **Media Tracking**: Track movies, books, anime, games, TV shows, and videos
- **Progress Management**: Set status (planned, in progress, completed, on hold, dropped)
- **Rating System**: Rate media with 1-10 scale
- **Review System**: Write detailed reviews with Markdown support
- **Collections**: Organize entries into custom collections
- **Search**: Real-time search across all media types

### 👤 User Experience
- **Guest Mode**: Use immediately without registration
- **Data Migration**: Seamlessly merge guest data into registered accounts
- **Public Sharing**: Share collections and profiles via public links
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Dark/Light Theme**: User preference support

### 🔧 Technical Features
- **Real-time Sync**: Automatic data synchronization
- **Offline Support**: Guest mode works without internet
- **Export/Import**: JSON data export and import
- **RESTful API**: Clean, documented API endpoints
- **JWT Authentication**: Secure token-based auth
- **Database Migrations**: Versioned schema updates

## 📚 API Documentation

### Authentication
- `POST /api/auth/login` - Login with email
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
- `POST /api/entries/sync` - Sync entries

### Collections
- `GET /api/collections` - List collections
- `POST /api/collections` - Create collection
- `GET /api/collections/:id` - Get collection
- `PATCH /api/collections/:id` - Update collection
- `DELETE /api/collections/:id` - Delete collection
- `POST /api/collections/:id/share` - Create share link

### Guest Mode
- `POST /api/guest/snapshot` - Create guest data snapshot
- `POST /api/guest/merge` - Merge guest data to account

### Public Sharing
- `GET /s/:token` - View public share

## 🔧 Configuration

### Default User
The application comes with a default user for testing:
- **Email**: `admin@example.com`
- **Name**: Admin User
- **ID**: `550e8400-e29b-41d4-a716-446655440000`

You can login with this user immediately after running the migrations.

### Environment Variables
Copy `back/env.example` to `back/.env` and configure:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=media_tracker

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24
```

## 🐳 Docker Commands

### Production
```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild
docker-compose up -d --build
```

### Development
```bash
# Start with hot reload
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f backend

# Stop
docker-compose -f docker-compose.dev.yml down
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
cd back && go test -cover ./...

# Run specific test
cd back && go test ./internal/services
```

## 📦 Deployment

### Production Build
```bash
# Build production image
docker build -f back/Dockerfile -t media-tracker-backend ./back

# Run with production compose
docker-compose up -d
```

### Environment Variables for Production
- Set `JWT_SECRET` to a strong secret
- Configure database credentials
- Set up Redis password if needed
- Configure external API keys

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🆘 Troubleshooting

### Common Issues

**Database connection failed:**
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres
```

**Backend won't start:**
```bash
# Check if dependencies are available
docker-compose logs backend

# Rebuild the image
docker-compose build backend
```

**Hot reload not working:**
```bash
# Check Air configuration
docker-compose -f docker-compose.dev.yml logs backend

# Restart development environment
make dev-down && make dev-up
```
