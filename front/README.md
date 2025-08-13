# Media Tracker Frontend

A modern SvelteKit frontend for the Media Tracker application. Built with Svelte 4, TypeScript, and Tailwind CSS.

## Features

- **Modern UI**: Clean, responsive design with Tailwind CSS
- **Guest Mode**: Use without registration, data stored locally
- **Authentication**: Email-based login with JWT tokens
- **Real-time Search**: Search media items as you type
- **Status Management**: Track progress with visual status indicators
- **Markdown Reviews**: Write rich reviews with Markdown support
- **Responsive Design**: Works on desktop, tablet, and mobile
- **TypeScript**: Full type safety throughout the application

## Tech Stack

- **Framework**: SvelteKit 2.0
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Svelte stores
- **Markdown**: Marked.js
- **HTTP Client**: Fetch API
- **Build Tool**: Vite

## Quick Start

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The development server will start on `http://localhost:3000` with API proxy to the backend.

## Project Structure

```
front/
├── src/
│   ├── app.css              # Global styles with Tailwind
│   ├── app.html             # HTML template
│   ├── lib/
│   │   ├── components/      # Reusable components
│   │   │   ├── Header.svelte
│   │   │   ├── EntryCard.svelte
│   │   │   ├── AddEntryDialog.svelte
│   │   │   └── StatusTabs.svelte
│   │   ├── stores/          # Svelte stores
│   │   │   ├── auth.ts      # Authentication state
│   │   │   └── entries.ts   # Entries state
│   │   ├── types/           # TypeScript types
│   │   │   └── index.ts     # Shared types
│   │   └── utils/           # Utility functions
│   │       ├── api.ts       # API client
│   │       └── storage.ts   # Local storage utilities
│   └── routes/              # SvelteKit routes
│       ├── +layout.svelte   # Root layout
│       └── +page.svelte     # Dashboard page
├── static/                  # Static assets
├── package.json
├── svelte.config.js
├── vite.config.js
├── tailwind.config.js
└── tsconfig.json
```

## Key Components

### Header.svelte
Navigation bar with authentication and guest mode switching.

### EntryCard.svelte
Displays individual media entries with status, rating, and review preview.

### AddEntryDialog.svelte
Modal for adding new entries with media search and type selection.

### StatusTabs.svelte
Filter tabs for different entry statuses (planned, in progress, completed, etc.).

## State Management

### Auth Store
Manages authentication state, user data, and guest mode.

```typescript
import { auth } from '$stores/auth';

// Check if user is authenticated
$: isAuthenticated = $auth.isAuthenticated;

// Login
auth.login(user, token);

// Logout
auth.logout();
```

### Entries Store
Manages media entries with filtering and CRUD operations.

```typescript
import { entries } from '$stores/entries';

// Get entries by status
const plannedEntries = entries.getByStatus('planned');

// Add new entry
entries.addEntry(newEntry);
```

## API Integration

The frontend communicates with the Go backend through the API client in `src/lib/utils/api.ts`.

### Authentication
- `POST /api/auth/login` - Email-based login
- `GET /api/auth/me` - Get user profile

### Media
- `POST /api/media` - Create media item
- `GET /api/media/search` - Search media items

### Entries
- `GET /api/entries` - List user entries
- `POST /api/entries` - Create entry
- `PATCH /api/entries/:id` - Update entry
- `DELETE /api/entries/:id` - Delete entry

## Guest Mode

The application supports a guest mode where users can:

1. **Use without registration** - All data stored in localStorage
2. **Add entries locally** - No backend communication required
3. **Export/Import data** - JSON export/import functionality
4. **Switch to account** - Merge guest data into registered account

### Local Storage Structure

```typescript
interface GuestData {
  guestId: string;
  entries: Entry[];
  media: MediaItem[];
}
```

## Development

### Available Scripts

```bash
npm run dev          # Start development server
npm run build        # Build for production
npm run preview      # Preview production build
npm run check        # Type check
npm run lint         # Lint code
npm run format       # Format code
```

### Environment Variables

Create a `.env` file in the frontend directory:

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080
```

### API Proxy

The development server proxies API requests to the backend:

```javascript
// vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

## Styling

The application uses Tailwind CSS with custom components:

```css
@layer components {
  .btn {
    @apply px-4 py-2 rounded-lg font-medium transition-colors duration-200;
  }
  
  .btn-primary {
    @apply bg-primary-600 text-white hover:bg-primary-700;
  }
  
  .card {
    @apply bg-white rounded-lg shadow-md p-6;
  }
}
```

## TypeScript

Full TypeScript support with shared types matching the backend:

```typescript
export type MediaType = 'video' | 'book' | 'anime' | 'game' | 'tv' | 'movie';
export type Status = 'planned' | 'in_progress' | 'completed' | 'on_hold' | 'dropped';
```

## Deployment

### Build for Production

```bash
npm run build
```

The built application will be in the `build` directory.

### Docker

```dockerfile
FROM node:18-alpine

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "run", "preview"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details
