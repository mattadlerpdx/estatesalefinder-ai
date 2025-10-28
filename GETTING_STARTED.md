# Getting Started with EstateSaleFinder.ai

## 🚀 Quick Start Guide

This guide will help you get the estate sale marketplace running locally.

### Prerequisites

- **Node.js** 18+ (for frontend)
- **Go** 1.22+ (for backend)
- **Docker & Docker Compose** (recommended for easiest setup)
- **Firebase Account** (for authentication)
- **PostgreSQL** 15+ (if not using Docker)

---

## Option 1: Docker Compose (Recommended)

The fastest way to get started is using Docker Compose, which sets up everything automatically.

### Step 1: Clone & Navigate

```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
```

### Step 2: Set Up Firebase Credentials

You'll use the existing Firebase credentials from Cadence:

```bash
# Backend already has credentials/firebase-adminsdk.json from Cadence
# This is already copied and ready to use!
```

### Step 3: Configure Frontend Environment

```bash
cd frontend
cp .env.local.example .env.local
```

Edit `.env.local` and add your Firebase web app config:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080

# Get these from Firebase Console > Project Settings > Your apps > Web app
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=cadencescm.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=cadencescm
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=cadencescm.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your-sender-id
NEXT_PUBLIC_FIREBASE_APP_ID=your-app-id
```

### Step 4: Install Frontend Dependencies

```bash
npm install
cd ..
```

### Step 5: Start Everything!

```bash
make dev
```

This single command will:
- Start PostgreSQL database
- Run database migrations
- Start Go backend API on port 8080
- Start Next.js frontend on port 3000

### Step 6: Open Your Browser

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/health

---

## Option 2: Manual Setup (Without Docker)

### Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Set up environment
cp .env.example .env
# Edit .env with your local PostgreSQL credentials

# Start PostgreSQL (if not using Docker)
# Then run migrations
psql -U postgres -d estatesale_db -f migrations/001_initial_schema.sql

# Run the API
go run cmd/api/main.go
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Set up environment
cp .env.local.example .env.local
# Edit .env.local with Firebase config

# Run development server
npm run dev
```

---

## 🧪 Testing the API

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"healthy"}
```

### Get All Sales (Public)

```bash
curl http://localhost:8080/api/sales
```

### Create a Sale (Requires Auth)

First, register a user on the frontend at http://localhost:3000/auth/register

Then use the Firebase token to create a sale:

```bash
curl -X POST http://localhost:8080/api/sales/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_FIREBASE_TOKEN" \
  -d '{
    "title": "Estate Sale - Vintage Furniture",
    "description": "Large estate sale featuring mid-century modern furniture",
    "city": "Portland",
    "state": "OR",
    "zip_code": "97201",
    "start_date": "2025-02-01T09:00:00Z",
    "end_date": "2025-02-02T17:00:00Z",
    "sale_type": "estate_sale"
  }'
```

---

## 📁 Project Structure

```
estatesalefinder-ai/
├── backend/                    # Go API
│   ├── cmd/api/main.go         # Entry point
│   ├── internal/domain/sale/   # Sale domain logic
│   ├── internal/infrastructure/
│   └── migrations/             # Database schema
│
├── frontend/                   # Next.js app
│   ├── app/                    # App Router pages
│   │   ├── page.tsx            # Homepage
│   │   ├── sales/              # Sales listing
│   │   └── dashboard/          # Seller dashboard
│   ├── components/             # Reusable components
│   └── lib/                    # Utilities & Firebase
│
├── docker-compose.yml          # Local dev environment
├── Makefile                    # Dev commands
└── README.md                   # Full documentation
```

---

## 🛠️ Common Make Commands

```bash
make dev              # Start dev environment (rebuild)
make dev-no-build     # Start dev environment (no rebuild)
make down             # Stop all services
make logs             # View all logs
make logs-backend     # View backend logs
make logs-frontend    # View frontend logs
make db-shell         # Access PostgreSQL shell
make db-reset         # Reset database (deletes all data!)
make rebuild-frontend # Rebuild frontend only
make rebuild-backend  # Rebuild backend only
```

---

## 🔥 Firebase Setup

### Get Firebase Web Config

1. Go to [Firebase Console](https://console.firebase.google.com)
2. Select your project (cadencescm)
3. Click ⚙️ Settings > Project Settings
4. Scroll down to "Your apps"
5. If you don't have a web app, click "Add app" > Web
6. Copy the `firebaseConfig` values to `frontend/.env.local`

### Firebase Authentication

The backend already has Firebase Admin SDK configured (using Cadence's credentials).

To enable authentication:
1. Firebase Console > Authentication > Sign-in method
2. Enable Email/Password and/or Google sign-in

---

## 🗄️ Database

### Access Database Shell

```bash
make db-shell
```

Or manually:
```bash
docker compose exec postgres psql -U postgres -d estatesale_db
```

### Useful SQL Queries

```sql
-- View all sales
SELECT id, title, city, state, status FROM estate_sales;

-- View all users
SELECT id, email, user_type FROM users;

-- View sales with images
SELECT s.id, s.title, COUNT(i.id) as image_count
FROM estate_sales s
LEFT JOIN sale_images i ON s.id = i.sale_id
GROUP BY s.id, s.title;
```

### Reset Database

```bash
make db-reset
```

**WARNING**: This deletes ALL data!

---

## 🚢 Deployment

### Deploy to Google Cloud Run

```bash
# Deploy both backend and frontend
make deploy

# Or deploy individually
make deploy-backend
make deploy-frontend
```

### Before First Deployment:

1. **Create Cloud SQL Instance:**
   ```bash
   gcloud sql instances create estatesale-postgres-instance \
     --database-version=POSTGRES_15 \
     --tier=db-f1-micro \
     --region=us-west1
   ```

2. **Create Database:**
   ```bash
   gcloud sql databases create estatesale_db \
     --instance=estatesale-postgres-instance
   ```

3. **Update `backend/env_vars.yaml`** with production credentials

4. **Deploy:**
   ```bash
   make deploy
   ```

---

## 🐛 Troubleshooting

### Frontend not updating?

```bash
make rebuild-frontend
# Then hard refresh browser: Ctrl+Shift+R
```

### Backend not starting?

Check logs:
```bash
make logs-backend
```

Common issues:
- Database not ready (wait 10 seconds after `make dev`)
- Firebase credentials missing
- Port 8080 already in use

### Database connection errors?

```bash
# Check if PostgreSQL is running
docker compose ps

# Restart database
docker compose restart postgres
```

### "Module not found" errors in frontend?

```bash
cd frontend
npm install
cd ..
make rebuild-frontend
```

---

## 📚 Next Steps

### MVP Features to Build:

1. **Sales Listing Page** (`/sales`)
   - Browse all published sales
   - Filter by city, state, date
   - Pagination

2. **Sale Detail Page** (`/sales/[id]`)
   - View sale details
   - Image gallery
   - Map with address
   - Save to favorites (auth required)

3. **Seller Dashboard** (`/dashboard`)
   - Create new sale
   - Manage existing sales
   - Upload images
   - View analytics

4. **Authentication** (`/auth/login`, `/auth/register`)
   - Email/password signup
   - Google sign-in
   - Protected routes

### Refer to ROADMAP.md for full development plan!

---

## 📞 Support

- **Documentation**: See README.md and ROADMAP.md
- **Issues**: Create issues in your GitHub repo
- **Firebase**: Check Cadence project for existing setup

---

**Built with ❤️ to modernize the estate sale industry**
