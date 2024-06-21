package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"golang.org/x/time/rate"

	gowebly "github.com/gowebly/helpers"
)

//go:embed all:static
var static embed.FS

var limiter = rate.NewLimiter(1, 5)

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7000"))
	if err != nil {
		return err
	}

	// Handle static files from the embed FS (with a custom handler).
	http.Handle("GET /static/", gowebly.StaticFileServerHandler(http.FS(static)))

	// Handle index page view.
	http.HandleFunc("GET /", indexViewHandler)

	http.HandleFunc("POST /api/search", rateLimit(fetchGoogleResults))

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
            return
        }
        next(w, r)
    }
}
