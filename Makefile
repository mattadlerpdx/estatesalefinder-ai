.PHONY: help dev dev-no-build up down restart logs build rebuild-frontend rebuild-backend deploy deploy-all deploy-backend deploy-frontend migrate clean

# Default target
help:
	@echo "EstateSaleFinder.ai Development Commands:"
	@echo ""
	@echo "Local Development:"
	@echo "  make dev              - Rebuild and start development environment"
	@echo "  make dev-no-build     - Start without rebuilding (faster, use if no code changes)"
	@echo "  make up               - Start services in background (with rebuild)"
	@echo "  make down             - Stop all services"
	@echo "  make restart          - Restart all services (with rebuild)"
	@echo "  make logs             - View logs from all services"
	@echo "  make logs-backend     - View backend logs"
	@echo "  make logs-frontend    - View frontend logs"
	@echo "  make build            - Rebuild all containers"
	@echo "  make rebuild-frontend - Rebuild only frontend (faster)"
	@echo "  make rebuild-backend  - Rebuild only backend (faster)"
	@echo "  make clean            - Stop services and remove volumes"
	@echo ""
	@echo "Database:"
	@echo "  make db-shell         - Access PostgreSQL shell"
	@echo "  make db-reset         - Reset database (WARNING: deletes data)"
	@echo "  make migrate          - Run database migrations"
	@echo ""
	@echo "Deployment (Google Cloud Run):"
	@echo "  make deploy                - Deploy both backend and frontend (recommended)"
	@echo "  make deploy-all            - Same as 'make deploy'"
	@echo "  make deploy-backend        - Build and deploy backend only (Cloud Build)"
	@echo "  make deploy-frontend       - Build and deploy frontend only"
	@echo "  make deploy-backend-manual - Build backend locally, then push and deploy"
	@echo ""
	@echo "Cloud Logs:"
	@echo "  make cloud-logs-backend   - View backend logs from Cloud Run"
	@echo "  make cloud-logs-frontend  - View frontend logs from Cloud Run"

# ==========================================
# LOCAL DEVELOPMENT
# ==========================================

dev:
	@echo "Rebuilding and starting local development environment..."
	docker compose up --build

dev-no-build:
	@echo "Starting local development environment (no rebuild)..."
	docker compose up

up:
	@echo "Starting services in background (with rebuild)..."
	docker compose up --build -d

down:
	@echo "Stopping all services..."
	docker compose down

restart:
	@echo "Restarting services (with rebuild)..."
	docker compose down
	docker compose up --build -d
	@echo "Services restarted! View logs with: make logs"

logs:
	docker compose logs -f

logs-backend:
	docker compose logs -f backend

logs-frontend:
	docker compose logs -f frontend

build:
	@echo "Rebuilding containers..."
	docker compose build

rebuild-frontend:
	@echo "Rebuilding frontend only..."
	docker compose up --build -d frontend
	@echo "Frontend rebuilt! View logs with: make logs-frontend"

rebuild-backend:
	@echo "Rebuilding backend only..."
	docker compose up --build -d backend
	@echo "Backend rebuilt! View logs with: make logs-backend"

clean:
	@echo "Stopping services and removing volumes..."
	docker compose down -v

# ==========================================
# DATABASE
# ==========================================

db-shell:
	docker compose exec postgres psql -U postgres -d estatesale_db

db-reset:
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose down -v; \
		docker compose up -d postgres; \
		echo "Database reset complete"; \
	fi

migrate:
	@echo "Running database migrations..."
	docker compose exec postgres psql -U postgres -d estatesale_db -f /docker-entrypoint-initdb.d/001_initial_schema.sql
	@echo "Migrations complete!"

# ==========================================
# DEPLOYMENT - BACKEND
# ==========================================

# Deploy backend using Cloud Build (recommended - faster, no local Docker needed)
deploy-backend:
	@echo "Building and pushing backend using Cloud Build..."
	gcloud builds submit --tag gcr.io/cadencescm/estatesale-backend-image backend/
	@echo "Deploying to Cloud Run..."
	gcloud run deploy estatesale-backend-service \
		--image gcr.io/cadencescm/estatesale-backend-image:latest \
		--platform managed \
		--region us-west1 \
		--allow-unauthenticated \
		--add-cloudsql-instances cadencescm:us-west1:estatesale-postgres-instance \
		--env-vars-file backend/env_vars.yaml
	@echo "Backend deployment complete!"

# Deploy backend manually (build locally, push to GCR, deploy)
deploy-backend-manual:
	@echo "Step 1: Building backend Docker image locally..."
	cd backend && docker build -t gcr.io/cadencescm/estatesale-backend-image .
	@echo "Step 2: Pushing to Google Container Registry..."
	docker push gcr.io/cadencescm/estatesale-backend-image:latest
	@echo "Step 3: Deploying to Cloud Run..."
	gcloud run deploy estatesale-backend-service \
		--image gcr.io/cadencescm/estatesale-backend-image:latest \
		--platform managed \
		--region us-west1 \
		--allow-unauthenticated \
		--add-cloudsql-instances cadencescm:us-west1:estatesale-postgres-instance \
		--env-vars-file backend/env_vars.yaml
	@echo "Backend deployment complete!"

# ==========================================
# DEPLOYMENT - FRONTEND
# ==========================================

deploy-frontend:
	@echo "Building frontend Docker image..."
	cd frontend && docker build -t gcr.io/cadencescm/estatesale-frontend-image . --no-cache
	@echo "Pushing to Google Container Registry..."
	docker push gcr.io/cadencescm/estatesale-frontend-image:latest
	@echo "Deploying to Cloud Run..."
	gcloud run deploy estatesale-frontend-service \
		--image gcr.io/cadencescm/estatesale-frontend-image:latest \
		--platform managed \
		--region us-west1 \
		--allow-unauthenticated \
		--port 3000 \
		--env-vars-file frontend/env_vars.yaml
	@echo "Frontend deployment complete!"

# ==========================================
# DEPLOYMENT - BOTH
# ==========================================

deploy-all: deploy-backend deploy-frontend
	@echo "=========================================="
	@echo "Full deployment complete!"
	@echo "Backend: https://estatesale-backend-service-[PROJECT-ID].us-west1.run.app"
	@echo "Frontend: https://estatesale-frontend-service-[PROJECT-ID].us-west1.run.app"
	@echo "=========================================="

# Shorthand alias for deploy-all
deploy: deploy-all

# ==========================================
# CLOUD LOGS
# ==========================================

cloud-logs-backend:
	gcloud run services logs read estatesale-backend-service \
		--region us-west1 \
		--limit 50

cloud-logs-frontend:
	gcloud run services logs read estatesale-frontend-service \
		--region us-west1 \
		--limit 50
