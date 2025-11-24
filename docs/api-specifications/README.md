# API Specifications

This directory contains comprehensive API specifications for all microservices in the Tabelogo v2 platform.

## Purpose

- **Single Source of Truth**: Centralized documentation for all API endpoints
- **Frontend Integration**: Clear contracts for frontend developers
- **Testing Reference**: Specifications for API testing
- **Onboarding**: Quick reference for new team members

## Structure

Each service has its own specification file:

```
docs/api-specifications/
├── README.md                 # This file
├── auth-service.md          # Authentication & user management
├── map-service.md           # Google Maps integration (TODO)
├── restaurant-service.md    # Restaurant data & favorites (TODO)
└── booking-service.md       # Reservation management (TODO)
```

## Specification Format

Each specification includes:

1. **Base URL**: Development and production endpoints
2. **Authentication**: Required headers and tokens
3. **Endpoints**: Complete list with:
   - HTTP method and path
   - Request body schema
   - Success response schema
   - Error responses with codes
4. **Data Models**: TypeScript-style type definitions
5. **Error Codes**: Standardized error codes
6. **Integration Notes**: Frontend-specific guidance

## Available Specifications

### ✅ Auth Service
- **File**: [auth-service.md](./auth-service.md)
- **Status**: Complete
- **Endpoints**: 4 (Register, Login, Refresh, Validate)
- **Port**: 8081

### 🔄 Map Service (TODO)
- **Port**: 8080
- **Endpoints**: Quick Search, Advance Search

### 🔄 Restaurant Service (TODO)
- **Port**: 8082
- **Endpoints**: Get Restaurant, Favorites CRUD

### 🔄 Booking Service (TODO)
- **Port**: 8083
- **Endpoints**: Create Booking, Get Bookings, Cancel

## Usage

### For Frontend Developers

1. Read the specification before implementing
2. Use the TypeScript types as reference
3. Implement error handling for all error codes
4. Test with actual backend service

### For Backend Developers

1. Keep specifications updated when changing APIs
2. Follow the documented response formats
3. Use consistent error codes
4. Add examples for complex requests

## Contributing

When adding or updating a specification:

1. Follow the existing format
2. Include request/response examples
3. Document all error cases
4. Add TypeScript type definitions
5. Include integration notes

## Related Documentation

- [Architecture](../../architecture.md) - Overall system design
- [Frontend README](../../web/README.md) - Frontend setup
- [Backend ENV Configuration](../../BACKEND_ENV_CONFIGURATION.md) - Environment setup
