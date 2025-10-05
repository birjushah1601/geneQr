# ABY-MED Admin Dashboard

## 🎯 Overview

Admin dashboard for managing manufacturer onboarding, equipment registry, field engineers, and service tickets.

## 🚀 Tech Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **UI Library:** shadcn/ui + Tailwind CSS
- **State Management:** React Query + Zustand
- **Forms:** React Hook Form + Zod validation
- **File Upload:** react-dropzone
- **Real-time:** Socket.io (for ticket notifications)
- **Authentication:** Next-Auth (ready for Keycloak)

## 📁 Project Structure

```
admin-ui/
├── src/
│   ├── app/                    # Next.js app router
│   │   ├── (auth)/
│   │   │   └── login/
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx      # Dashboard layout
│   │   │   ├── page.tsx        # Overview dashboard
│   │   │   ├── manufacturers/  # Manufacturer management
│   │   │   ├── equipment/      # Equipment registry
│   │   │   ├── engineers/      # Field engineers
│   │   │   ├── tickets/        # Service tickets
│   │   │   └── settings/       # Settings
│   │   └── api/                # API routes
│   ├── components/
│   │   ├── ui/                 # shadcn/ui components
│   │   ├── forms/              # Form components
│   │   ├── tables/             # Data tables
│   │   └── dashboard/          # Dashboard widgets
│   ├── lib/
│   │   ├── api/                # API client
│   │   ├── hooks/              # Custom React hooks
│   │   ├── utils/              # Utilities
│   │   └── validation/         # Zod schemas
│   ├── types/                  # TypeScript types
│   └── styles/                 # Global styles
├── public/
└── package.json
```

## 🔧 Setup Instructions

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
KEYCLOAK_CLIENT_ID=aby-med-admin
KEYCLOAK_CLIENT_SECRET=your-client-secret
KEYCLOAK_ISSUER=http://localhost:8080/realms/aby-med
```

### 3. Run Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

## 📋 Features

### Phase 1 (Current)
- ✅ Manufacturer onboarding with CSV upload
- ✅ Equipment registry management
- ✅ Field engineer management
- ✅ Service ticket dashboard
- ✅ Manual engineer assignment
- ✅ Service overview

### Phase 2 (Next)
- 🔄 WhatsApp integration
- 🔄 Real-time ticket updates
- 🔄 Advanced filtering
- 🔄 Reporting dashboard

### Phase 3 (Future)
- ⏳ Keycloak integration
- ⏳ Role-based access control
- ⏳ Mobile responsive views
- ⏳ Engineer mobile app

## 📱 Screenshots

(Will be added after implementation)

## 🔗 API Integration

See `docs/API_INTEGRATION.md` for detailed API documentation.

## 🧪 Testing

```bash
# Run tests
npm test

# Run E2E tests
npm run test:e2e

# Type checking
npm run type-check
```

## 📦 Build & Deploy

```bash
# Build for production
npm run build

# Start production server
npm start
```

## 📄 License

Private - ABY-MED Platform
