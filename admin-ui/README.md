# ServQR Admin Dashboard

## ðŸŽ¯ Overview

Admin dashboard for managing manufacturer onboarding, equipment registry, field engineers, and service tickets.

## ðŸš€ Tech Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **UI Library:** shadcn/ui + Tailwind CSS
- **State Management:** React Query + Zustand
- **Forms:** React Hook Form + Zod validation
- **File Upload:** react-dropzone
- **Real-time:** Socket.io (for ticket notifications)
- **Authentication:** Next-Auth (ready for Keycloak)

## ðŸ“ Project Structure

```
admin-ui/
â”œâ”€â”€ src/
â”‚   â”œâ”€â”€ app/                    # Next.js app router
â”‚   â”‚   â”œâ”€â”€ (auth)/
â”‚   â”‚   â”‚   â””â”€â”€ login/
â”‚   â”‚   â”œâ”€â”€ (dashboard)/
â”‚   â”‚   â”‚   â”œâ”€â”€ layout.tsx      # Dashboard layout
â”‚   â”‚   â”‚   â”œâ”€â”€ page.tsx        # Overview dashboard
â”‚   â”‚   â”‚   â”œâ”€â”€ manufacturers/  # Manufacturer management
â”‚   â”‚   â”‚   â”œâ”€â”€ equipment/      # Equipment registry
â”‚   â”‚   â”‚   â”œâ”€â”€ engineers/      # Field engineers
â”‚   â”‚   â”‚   â”œâ”€â”€ tickets/        # Service tickets
â”‚   â”‚   â”‚   â””â”€â”€ settings/       # Settings
â”‚   â”‚   â””â”€â”€ api/                # API routes
â”‚   â”œâ”€â”€ components/
â”‚   â”‚   â”œâ”€â”€ ui/                 # shadcn/ui components
â”‚   â”‚   â”œâ”€â”€ forms/              # Form components
â”‚   â”‚   â”œâ”€â”€ tables/             # Data tables
â”‚   â”‚   â””â”€â”€ dashboard/          # Dashboard widgets
â”‚   â”œâ”€â”€ lib/
â”‚   â”‚   â”œâ”€â”€ api/                # API client
â”‚   â”‚   â”œâ”€â”€ hooks/              # Custom React hooks
â”‚   â”‚   â”œâ”€â”€ utils/              # Utilities
â”‚   â”‚   â””â”€â”€ validation/         # Zod schemas
â”‚   â”œâ”€â”€ types/                  # TypeScript types
â”‚   â””â”€â”€ styles/                 # Global styles
â”œâ”€â”€ public/
â””â”€â”€ package.json
```

## ðŸ”§ Setup Instructions

### 1. Install Dependencies

```bash
cd admin-ui
npm install
```

### 2. Environment Variables

Create `.env.local`:

```env
# Backend API
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081
NEXT_PUBLIC_WS_URL=ws://localhost:8081

# Authentication (for later)
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key

# Keycloak (for later)
KEYCLOAK_CLIENT_ID=ServQR-admin
KEYCLOAK_CLIENT_SECRET=your-client-secret
KEYCLOAK_ISSUER=http://localhost:8080/realms/ServQR
```

### 3. Run Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

## ðŸ“‹ Features

### Phase 1 (Current)
- âœ… Manufacturer onboarding with CSV upload
- âœ… Equipment registry management
- âœ… Field engineer management
- âœ… Service ticket dashboard
- âœ… Manual engineer assignment
- âœ… Service overview

### Phase 2 (Next)
- ðŸ”„ WhatsApp integration
- ðŸ”„ Real-time ticket updates
- ðŸ”„ Advanced filtering
- ðŸ”„ Reporting dashboard

### Phase 3 (Future)
- â³ Keycloak integration
- â³ Role-based access control
- â³ Mobile responsive views
- â³ Engineer mobile app

## ðŸ“± Screenshots

(Will be added after implementation)

## ðŸ”— API Integration

See `docs/API_INTEGRATION.md` for detailed API documentation.

## ðŸ§ª Testing

```bash
# Run tests
npm test

# Run E2E tests
npm run test:e2e

# Type checking
npm run type-check
```

## ðŸ“¦ Build & Deploy

```bash
# Build for production
npm run build

# Start production server
npm start
```

## ðŸ“„ License

Private - ServQR Platform
