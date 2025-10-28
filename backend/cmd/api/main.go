package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/listing"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/user"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/cache"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/controllers"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/db/postgres"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/middleware"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/scraper"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Initialize Firebase
	if err := middleware.InitFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Retrieve environment variables for database connection
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST") // Can be Unix socket path or IP address
	dbPort := os.Getenv("DB_PORT") // Port number (default 5432 for PostgreSQL)

	// Ensure all required environment variables are set
	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		log.Fatal("Database environment variables are not set properly.")
	}

	// Default port if not specified
	if dbPort == "" {
		dbPort = "5432"
	}

	// Build DSN - use port if it's a TCP connection (IP address), otherwise Unix socket
	var dsn string
	if dbPort != "" && !strings.HasPrefix(dbHost, "/cloudsql") {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	} else {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPass, dbName, dbHost)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize repositories
	listingRepo := postgres.NewListingRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	listingService := listing.NewService(listingRepo)
	userService := user.NewService(userRepo)

	// Initialize Redis cache
	redisClient := cache.NewRedisClient()
	defer redisClient.Close()

	// Initialize scraper service (with repository for hybrid storage)
	scraperService := scraper.NewScraperService(redisClient, listingRepo)

	// Initialize handlers
	listingHandler := controllers.NewListingHandler(listingService, userService)
	listingHandler.SetScraperService(scraperService)
	userHandler := controllers.NewUserHandler(userService)

	// Set up the router using stdlib http.ServeMux
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Helper function to apply auth middleware to handlers
	authMiddleware := func(handler http.HandlerFunc) http.Handler {
		return middleware.FirebaseMiddleware(http.HandlerFunc(handler))
	}

	// ==========================================
	// PUBLIC ENDPOINTS (No auth required)
	// ==========================================

	// Sales - Public browsing (owned + scraped)
	mux.Handle("/api/sales", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listingHandler.GetAggregatedSales(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Individual sale - Public viewing
	mux.Handle("/api/sales/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a specific sale ID (not /api/sales/something-else)
		path := r.URL.Path
		if strings.HasPrefix(path, "/api/sales/") && !strings.Contains(strings.TrimPrefix(path, "/api/sales/"), "/") {
			if r.Method == http.MethodGet {
				listingHandler.GetByID(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})))

	// ==========================================
	// AUTHENTICATED ENDPOINTS
	// ==========================================

	// User management
	mux.Handle("/ensureUser", corsMiddleware(authMiddleware(userHandler.EnsureUser)))

	// Sales - Create (authenticated sellers only)
	mux.Handle("/api/sales/create", corsMiddleware(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			listingHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// My Sales - Get user's own sales
	mux.Handle("/api/my-sales", corsMiddleware(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listingHandler.GetMySales(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Update sale
	mux.Handle("/api/sales/update/", corsMiddleware(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			listingHandler.Update(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Delete sale
	mux.Handle("/api/sales/delete/", corsMiddleware(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			listingHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Add image to sale
	mux.Handle("/api/sales/images/", corsMiddleware(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			listingHandler.AddImage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Get PORT from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGIN")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000"
		}

		// Check if request origin is in allowed list
		origin := r.Header.Get("Origin")
		if origin != "" {
			for _, allowed := range strings.Split(allowedOrigins, ",") {
				if strings.TrimSpace(allowed) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
