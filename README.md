# EstateSaleFinder.ai

> AI-Powered Estate Sale Marketplace - Connecting Buyers and Sellers Nationwide

## Overview

EstateSaleFinder.ai is a modern, full-stack estate sale marketplace platform designed to compete with and improve upon legacy platforms like estatesale-finder.com. Built with cutting-edge technology and AI-powered features, we offer a superior user experience at significantly lower prices.

### Competitive Advantages

- **Major Price Undercut**: Subscription plans 40-60% cheaper than competitors
- **AI-Powered Search**: Smart recommendations and intelligent item categorization
- **National Coverage**: Not limited to specific regions (vs competitor's Pacific Northwest focus)
- **Modern UX**: Built with Next.js 14 and Tailwind CSS for lightning-fast performance
- **Mobile-First**: Responsive design optimized for all devices
- **SEO Optimized**: Server-side rendering for maximum discoverability

## Tech Stack

### Frontend
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript (optional) / JavaScript
- **Styling**: Tailwind CSS 3.4+
- **UI Components**: Headless UI, Heroicons
- **State Management**: React Context API
- **Data Fetching**: SWR / TanStack Query
- **Forms**: React Hook Form
- **Authentication**: Firebase Auth

### Backend
- **Language**: Go 1.22+
- **Architecture**: Clean Architecture + Domain-Driven Design (DDD)
- **Database**: PostgreSQL 15
- **Auth**: Firebase Admin SDK
- **API**: RESTful with Gorilla Mux
- **File Storage**: AWS S3 / Google Cloud Storage
- **Payment Processing**: Stripe

### Infrastructure
- **Frontend Hosting**: Vercel / Netlify
- **Backend Hosting**: Google Cloud Run / AWS ECS
- **Database**: Managed PostgreSQL (Cloud SQL / RDS)
- **CDN**: Cloudflare / CloudFront
- **CI/CD**: GitHub Actions
- **Monitoring**: Sentry, Google Analytics

## Project Structure

```
estatesalefinder-ai/
├── backend/                    # Go API server
│   ├── cmd/api/                # Application entry point
│   ├── internal/
│   │   ├── domain/             # Business logic (DDD)
│   │   │   ├── user/
│   │   │   ├── sale/
│   │   │   ├── item/
│   │   │   ├── professional/
│   │   │   ├── review/
│   │   │   └── subscription/
│   │   └── infrastructure/     # External concerns
│   │       ├── controllers/    # HTTP handlers
│   │       ├── db/postgres/    # Repository implementations
│   │       ├── middleware/     # Auth, CORS, logging
│   │       ├── storage/        # Image upload service
│   │       └── payment/        # Stripe integration
│   ├── migrations/             # Database migrations
│   └── Dockerfile
│
├── frontend/                   # Next.js application
│   ├── app/                    # Next.js 14 App Router
│   │   ├── (public)/           # Public routes
│   │   │   ├── page.tsx        # Homepage
│   │   │   ├── sales/          # Browse sales
│   │   │   ├── professionals/  # Professional directory
│   │   │   └── about/
│   │   ├── (auth)/             # Authenticated routes
│   │   │   ├── dashboard/      # Seller dashboard
│   │   │   ├── my-sales/       # Manage listings
│   │   │   └── saved/          # Saved favorites
│   │   └── api/                # API routes (optional)
│   ├── components/
│   │   ├── shared/             # Reusable components
│   │   └── domain/             # Domain-specific components
│   ├── contexts/               # React contexts
│   ├── services/               # API integration
│   ├── lib/                    # Utilities
│   └── public/                 # Static assets
│
├── docs/                       # Additional documentation
├── .github/workflows/          # CI/CD pipelines
├── docker-compose.yml          # Local development
├── README.md                   # This file
└── ROADMAP.md                  # Development roadmap
```

## Core Features

### For Buyers (Public)
- **Browse Sales**: Search and filter estate sales by location, date, category
- **Map View**: Interactive map showing nearby sales
- **Sale Details**: Photos, descriptions, dates, directions
- **Save Favorites**: Bookmark interesting sales
- **Email Alerts**: Get notified of new sales in your area
- **Professional Directory**: Find auctioneers, appraisers, estate sale companies

### For Sellers (Authenticated)
- **Create Listings**: Post estate sales with rich descriptions and photos
- **Manage Sales**: Edit, cancel, or mark sales as completed
- **Analytics Dashboard**: Track views, saves, and engagement
- **Tiered Pricing**: Choose from Basic, Featured, or Premium listings
- **Subscription Plans**: Monthly/annual plans for frequent sellers
- **Photo Upload**: Drag-and-drop image uploads with automatic optimization

### For Professionals (Authenticated)
- **Business Profile**: Showcase services, reviews, service areas
- **Verified Badge**: Get verified for credibility
- **Review System**: Collect and display client testimonials
- **Lead Generation**: Connect with potential clients

### AI-Powered Features (Future)
- **Smart Categorization**: Auto-tag items using computer vision
- **Price Suggestions**: AI-estimated pricing based on market data
- **Personalized Recommendations**: Show relevant sales to buyers
- **Fraud Detection**: Identify suspicious listings

## Getting Started

### Prerequisites
- Node.js 18+ (for frontend)
- Go 1.22+ (for backend)
- PostgreSQL 15+
- Docker & Docker Compose (optional, for local dev)
- Firebase account (for authentication)

### Installation

#### 1. Clone the repository
```bash
git clone https://github.com/yourusername/estatesalefinder-ai.git
cd estatesalefinder-ai
```

#### 2. Backend Setup
```bash
cd backend

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
make migrate-up

# Start the API server
go run cmd/api/main.go
# Server runs on http://localhost:8080
```

#### 3. Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Set up environment variables
cp .env.example .env.local
# Edit .env.local with your configuration

# Start development server
npm run dev
# App runs on http://localhost:3000
```

#### 4. Using Docker Compose (Recommended)
```bash
# From project root
docker compose up --build

# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Database: localhost:5432
```

## Development Workflow

### Backend
```bash
# Run tests
go test ./...

# Run with hot reload (using air)
air

# Format code
go fmt ./...

# Build for production
go build -o api cmd/api/main.go
```

### Frontend
```bash
# Run development server
npm run dev

# Run type checking (if using TypeScript)
npm run type-check

# Build for production
npm run build

# Start production server
npm start
```

## Environment Variables

### Backend (.env)
```bash
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/estatesale_db
FIREBASE_CREDENTIALS_PATH=./credentials/firebase-adminsdk.json
CORS_ALLOWED_ORIGIN=http://localhost:3000
AWS_S3_BUCKET=estatesale-images
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
STRIPE_SECRET_KEY=sk_test_...
```

### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_FIREBASE_API_KEY=your-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-app.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...
```

## Deployment

### Frontend (Vercel)
```bash
# Connect your GitHub repo to Vercel
# Add environment variables in Vercel dashboard
# Deploy automatically on push to main branch
```

### Backend (Google Cloud Run)
```bash
# Build and push Docker image
docker build -t gcr.io/your-project/estatesale-api:latest ./backend
docker push gcr.io/your-project/estatesale-api:latest

# Deploy to Cloud Run
gcloud run deploy estatesale-api \
  --image gcr.io/your-project/estatesale-api:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## Database Schema

See full schema in `backend/migrations/001_initial_schema.sql`

**Core Tables:**
- `users` - User accounts (buyers, sellers, professionals)
- `user_profiles` - Extended user information
- `estate_sales` - Main listing entity
- `sale_images` - Photos for listings
- `sale_items` - Individual items within sales
- `saved_sales` - User favorites
- `professionals` - Business directory
- `reviews` - Professional reviews
- `subscription_plans` - Pricing tiers
- `user_subscriptions` - Active subscriptions

## API Endpoints

### Public Endpoints
- `GET /api/sales` - List all published sales (with filters)
- `GET /api/sales/:id` - Get sale details
- `GET /api/professionals` - List professionals
- `GET /api/professionals/:id` - Get professional profile
- `GET /api/health` - Health check

### Authenticated Endpoints (require Bearer token)
- `POST /api/sales` - Create new sale
- `PUT /api/sales/:id` - Update sale
- `DELETE /api/sales/:id` - Delete sale
- `POST /api/sales/:id/images` - Upload images
- `POST /api/saved-sales` - Save a sale
- `DELETE /api/saved-sales/:id` - Unsave a sale
- `GET /api/my-sales` - Get user's sales
- `POST /api/subscriptions` - Create subscription
- `POST /api/reviews` - Create review

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style
- **Backend**: Follow Go conventions, use `gofmt`
- **Frontend**: ESLint + Prettier configured
- **Commits**: Use conventional commits (feat:, fix:, docs:, etc.)

## License

MIT License - see LICENSE file for details

## Support

- **Documentation**: [docs/](./docs/)
- **Issues**: [GitHub Issues](https://github.com/yourusername/estatesalefinder-ai/issues)
- **Email**: support@estatesalefinder.ai

## Roadmap

See [ROADMAP.md](./ROADMAP.md) for detailed development phases and feature timeline.

## Acknowledgments

- Inspired by architecture patterns from the Cadence project
- Built with modern tools and frameworks from the open-source community

---

**Built with ❤️ to modernize the estate sale industry**
