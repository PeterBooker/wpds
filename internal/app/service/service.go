package http

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/gorilla/websocket"
)

var (
	t = template.Must(template.ParseFiles("../../frontend/templates/app.html"))

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Init registers routes and starts the HTTP server.
func Init(ctx context.Context, valv *valve.Valve) {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// File Server
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "../../frontend/build/static")
	FileServer(r, "/static", http.Dir(filesDir))

	r.Get("/", appHandler)
	r.Get("/search/", appHandler)
	r.Get("/search/{id}", appHandler)
	r.Get("/search/list", appHandler)
	r.Get("/about", appHandler)

	// Test API Endpoints
	r.Get("/api/page/dashboard/", apiPageDashboardHandler)
	r.Get("/api/page/search/list/", apiPageSearchListHandler)
	r.Get("/api/page/search/{id}/", apiPageSearchViewHandler)

	srv := &http.Server{
		Handler:      chi.ServerBaseContext(ctx, r),
		Addr:         ":50505",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {

		for range c {

			// sig is a ^C, handle it
			fmt.Println("Shutting Down...")

			// first valv
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// Start HTTP Shutdown
			srv.Shutdown(ctx)

			// Verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				fmt.Println("Not all connections closed.")
			case <-ctx.Done():

			}

		}

	}()

	log.Fatal(srv.ListenAndServe())

}
