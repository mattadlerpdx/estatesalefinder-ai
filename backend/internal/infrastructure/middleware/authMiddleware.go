package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

type ContextKey string

const ContextKeyUID ContextKey = "uid"

// InitFirebase sets up the Firebase Admin SDK.
// It checks APP_ENV to determine whether to use local credentials or default GCP credentials.
func InitFirebase() error {
	env := os.Getenv("APP_ENV") // "local" or "production"
	var app *firebase.App
	var err error

	if env == "production" {
		// Use default credentials on GCP (e.g., Cloud Run, App Engine, etc.)
		app, err = firebase.NewApp(context.Background(), nil)
	} else {
		// Local development: use JSON credentials file
		opt := option.WithCredentialsFile("credentials/firebase-adminsdk.json")
		app, err = firebase.NewApp(context.Background(), nil, opt)
	}

	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	firebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing firebase auth client: %v", err)
	}

	return nil
}

// FirebaseMiddleware validates the Bearer token and injects UID into the request context.
func FirebaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := firebaseAuth.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Set UID in context for downstream use
		ctx := context.WithValue(r.Context(), ContextKeyUID, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
