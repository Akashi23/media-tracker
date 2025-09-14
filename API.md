# Media Tracker API Documentation

## Overview

The Media Tracker API is a RESTful service built with Go and Gin framework. It provides endpoints for managing media consumption tracking, including authentication, media management, entry tracking, collections, and sharing.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://your-domain.com`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "data": <response_data>,
  "message": "Success message"
}
```

### Error Response
```json
{
  "error": "Error message",
  "details": "Additional error details"
}
```

## Endpoints

### Authentication

#### Login
```http
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "User Name",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "Login successful"
}
```

#### Logout
```http
POST /api/auth/logout
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "message": "Logout successful"
}
```

#### Get Profile
```http
GET /api/auth/me
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "User Name",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Media

#### Create Media Item
```http
POST /api/media
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Movie Title",
  "type": "movie",
  "description": "Movie description",
  "release_date": "2024-01-01",
  "genres": ["Action", "Drama"],
  "external_id": "tmdb_12345"
}
```

**Response:**
```json
{
  "data": {
    "id": "media-uuid",
    "title": "Movie Title",
    "type": "movie",
    "description": "Movie description",
    "release_date": "2024-01-01",
    "genres": ["Action", "Drama"],
    "external_id": "tmdb_12345",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Media created successfully"
}
```

#### Update Media Item
```http
PUT /api/media/:id
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:** Same as create media

**Response:** Same as create media

#### Search Media
```http
GET /api/media/search?q=query&type=movie
```

**Query Parameters:**
- `q` (string): Search query
- `type` (string, optional): Media type filter (movie, book, anime, game, tv, video)

**Response:**
```json
{
  "data": [
    {
      "id": "media-uuid",
      "title": "Movie Title",
      "type": "movie",
      "description": "Movie description",
      "release_date": "2024-01-01",
      "genres": ["Action", "Drama"],
      "external_id": "tmdb_12345"
    }
  ],
  "message": "Search completed"
}
```

### Entries

#### List Entries
```http
GET /api/entries
```

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `status` (string, optional): Filter by status (planned, in_progress, completed, on_hold, dropped)
- `type` (string, optional): Filter by media type
- `limit` (int, optional): Number of entries to return (default: 50)
- `offset` (int, optional): Number of entries to skip (default: 0)

**Response:**
```json
{
  "data": [
    {
      "id": "entry-uuid",
      "user_id": "user-uuid",
      "media_id": "media-uuid",
      "status": "completed",
      "rating": 8,
      "review": "Great movie!",
      "progress": 100,
      "started_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-02T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-02T00:00:00Z",
      "media": {
        "id": "media-uuid",
        "title": "Movie Title",
        "type": "movie",
        "description": "Movie description",
        "release_date": "2024-01-01",
        "genres": ["Action", "Drama"]
      }
    }
  ],
  "message": "Entries retrieved successfully"
}
```

#### Create Entry
```http
POST /api/entries
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "media_id": "media-uuid",
  "status": "planned",
  "rating": null,
  "review": "",
  "progress": 0
}
```

**Response:**
```json
{
  "data": {
    "id": "entry-uuid",
    "user_id": "user-uuid",
    "media_id": "media-uuid",
    "status": "planned",
    "rating": null,
    "review": "",
    "progress": 0,
    "started_at": null,
    "completed_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Entry created successfully"
}
```

#### Get Entry
```http
GET /api/entries/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response:** Same as create entry response

#### Update Entry
```http
PATCH /api/entries/:id
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "status": "completed",
  "rating": 8,
  "review": "Great movie!",
  "progress": 100
}
```

**Response:** Same as create entry response

#### Delete Entry
```http
DELETE /api/entries/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "message": "Entry deleted successfully"
}
```

#### Sync Entries
```http
POST /api/entries/sync
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "entries": [
    {
      "media_id": "media-uuid",
      "status": "completed",
      "rating": 8,
      "review": "Great movie!",
      "progress": 100
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "synced": 1,
    "created": 0,
    "updated": 1
  },
  "message": "Sync completed successfully"
}
```

### Collections

#### List Collections
```http
GET /api/collections
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "data": [
    {
      "id": "collection-uuid",
      "user_id": "user-uuid",
      "name": "My Favorites",
      "description": "My favorite movies and shows",
      "is_public": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "entry_count": 5
    }
  ],
  "message": "Collections retrieved successfully"
}
```

#### Create Collection
```http
POST /api/collections
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "My Favorites",
  "description": "My favorite movies and shows",
  "is_public": false,
  "entry_ids": ["entry-uuid-1", "entry-uuid-2"]
}
```

**Response:**
```json
{
  "data": {
    "id": "collection-uuid",
    "user_id": "user-uuid",
    "name": "My Favorites",
    "description": "My favorite movies and shows",
    "is_public": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "entries": [
      {
        "id": "entry-uuid-1",
        "status": "completed",
        "rating": 8,
        "review": "Great movie!",
        "progress": 100,
        "media": {
          "id": "media-uuid",
          "title": "Movie Title",
          "type": "movie",
          "description": "Movie description"
        }
      }
    ]
  },
  "message": "Collection created successfully"
}
```

#### Get Collection
```http
GET /api/collections/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response:** Same as create collection response

#### Update Collection
```http
PATCH /api/collections/:id
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "Updated Name",
  "description": "Updated description",
  "is_public": true,
  "entry_ids": ["entry-uuid-1", "entry-uuid-2", "entry-uuid-3"]
}
```

**Response:** Same as create collection response

#### Delete Collection
```http
DELETE /api/collections/:id
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "message": "Collection deleted successfully"
}
```

#### Create Share Link
```http
POST /api/collections/:id/share
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "data": {
    "token": "share-token-uuid",
    "expires_at": "2024-12-31T23:59:59Z",
    "url": "/s/share-token-uuid"
  },
  "message": "Share link created successfully"
}
```

### Guest Mode

#### Create Guest Snapshot
```http
POST /api/guest/snapshot
```

**Request Body:**
```json
{
  "guest_id": "guest-uuid",
  "entries": [
    {
      "media": {
        "title": "Movie Title",
        "type": "movie",
        "description": "Movie description"
      },
      "status": "completed",
      "rating": 8,
      "review": "Great movie!",
      "progress": 100
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "snapshot_id": "snapshot-uuid",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Guest snapshot created successfully"
}
```

#### Merge Guest Data
```http
POST /api/guest/merge
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "snapshot_id": "snapshot-uuid"
}
```

**Response:**
```json
{
  "data": {
    "merged_entries": 5,
    "merged_media": 3
  },
  "message": "Guest data merged successfully"
}
```

### Public Sharing

#### Get Public Share
```http
GET /s/:token
```

**Response:**
```json
{
  "data": {
    "collection": {
      "id": "collection-uuid",
      "name": "My Favorites",
      "description": "My favorite movies and shows",
      "created_at": "2024-01-01T00:00:00Z",
      "entries": [
        {
          "id": "entry-uuid",
          "status": "completed",
          "rating": 8,
          "review": "Great movie!",
          "progress": 100,
          "media": {
            "id": "media-uuid",
            "title": "Movie Title",
            "type": "movie",
            "description": "Movie description",
            "release_date": "2024-01-01",
            "genres": ["Action", "Drama"]
          }
        }
      ]
    },
    "share_info": {
      "token": "share-token-uuid",
      "expires_at": "2024-12-31T23:59:59Z"
    }
  }
}
```

## Data Types

### Media Types
- `movie` - Movies
- `book` - Books
- `anime` - Anime series
- `game` - Video games
- `tv` - TV shows
- `video` - Other videos

### Entry Status
- `planned` - Planned to watch/read/play
- `in_progress` - Currently watching/reading/playing
- `completed` - Finished
- `on_hold` - Paused
- `dropped` - Stopped

### Rating Scale
- 1-10 integer scale
- `null` for unrated entries

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation error |
| 500 | Internal Server Error - Server error |

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Authentication endpoints**: 5 requests per minute per IP
- **Other endpoints**: 100 requests per minute per authenticated user

## Pagination

List endpoints support pagination with the following parameters:
- `limit`: Number of items per page (default: 50, max: 100)
- `offset`: Number of items to skip (default: 0)

## Examples

### Complete Workflow

1. **Login**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

2. **Create Media**
```bash
curl -X POST http://localhost:8080/api/media \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "The Matrix",
    "type": "movie",
    "description": "A computer hacker learns about the true nature of reality.",
    "release_date": "1999-03-31",
    "genres": ["Action", "Sci-Fi"]
  }'
```

3. **Create Entry**
```bash
curl -X POST http://localhost:8080/api/entries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "media_id": "media-uuid",
    "status": "completed",
    "rating": 9,
    "review": "Mind-bending masterpiece!",
    "progress": 100
  }'
```

4. **Create Collection**
```bash
curl -X POST http://localhost:8080/api/collections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Sci-Fi Classics",
    "description": "My favorite science fiction movies",
    "is_public": true,
    "entry_ids": ["entry-uuid"]
  }'
```

5. **Share Collection**
```bash
curl -X POST http://localhost:8080/api/collections/collection-uuid/share \
  -H "Authorization: Bearer <token>"
```

## SDKs and Libraries

### JavaScript/TypeScript
```typescript
// Example API client
class MediaTrackerAPI {
  private baseURL: string;
  private token?: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  async login(email: string) {
    const response = await fetch(`${this.baseURL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email })
    });
    return response.json();
  }

  async getEntries() {
    const response = await fetch(`${this.baseURL}/api/entries`, {
      headers: { 'Authorization': `Bearer ${this.token}` }
    });
    return response.json();
  }
}
```

### Go
```go
// Example Go client
type MediaTrackerClient struct {
    baseURL string
    token   string
    client  *http.Client
}

func (c *MediaTrackerClient) Login(email string) (*LoginResponse, error) {
    data := map[string]string{"email": email}
    jsonData, _ := json.Marshal(data)
    
    req, _ := http.NewRequest("POST", c.baseURL+"/api/auth/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result LoginResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

## Changelog

### v1.0.0
- Initial API release
- Authentication with JWT
- Media management
- Entry tracking
- Collections
- Guest mode
- Public sharing
