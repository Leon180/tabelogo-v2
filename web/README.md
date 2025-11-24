# Tabelogo v2 Frontend

A modern, map-first restaurant discovery platform built with Next.js and Google Maps.

## Features

- **Map-First Interface**: Interactive Google Maps as the primary UI
- **Dual Search Modes**: 
  - Quick Search: Click markers on the map
  - Advance Search: Text search with filters (rating, open now, relevance/distance)
- **Real-time Filtering**: Location-based search using map bounds
- **Dark Mode Design**: Premium aesthetic with amber/gold accents
- **Responsive**: Works on desktop, tablet, and mobile

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: TailwindCSS + Shadcn/UI
- **Maps**: @vis.gl/react-google-maps
- **State Management**: React Query
- **API Client**: Axios

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Google Maps API Key (see ENV_CONFIG.md)

### Installation

```bash
# Install dependencies
npm install

# Create environment file
cp ENV_CONFIG.md .env.local
# Edit .env.local and add your Google Maps API key

# Run development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Environment Variables

See `ENV_CONFIG.md` for required environment variables.

## Project Structure

```
src/
├── app/                    # Next.js app router pages
│   ├── page.tsx           # Main map interface
│   ├── restaurants/       # Restaurant detail pages
│   ├── auth/              # Login/register pages
│   └── profile/           # User profile
├── components/
│   ├── Map/               # Google Maps components
│   ├── Search/            # Search form components
│   ├── Restaurant/        # Restaurant card components
│   └── ui/                # Shadcn UI components
├── lib/
│   └── api/               # API client modules
│       ├── map-service.ts
│       ├── auth-service.ts
│       ├── restaurant-service.ts
│       └── booking-service.ts
├── types/                 # TypeScript type definitions
├── hooks/                 # Custom React hooks
└── stores/                # State management
```

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

### API Integration

The frontend communicates with the following microservices:

- **Map Service** (port 8080): Google Maps API proxy
- **Auth Service** (port 8081): User authentication
- **Restaurant Service** (port 8082): Restaurant data
- **Booking Service** (port 8083): Reservation management

Configure service URLs in `.env.local`.

## Features Roadmap

- [x] Map-first interface
- [x] Advance search with filters
- [x] Dark mode design
- [ ] User authentication
- [ ] Restaurant detail pages
- [ ] Booking functionality
- [ ] Favorites management
- [ ] User profile
- [ ] Bilingual support (EN/JP)
- [ ] Mobile optimization

## License

Private project
