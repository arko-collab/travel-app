package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/travel-api/build/internal/approval"
	"github.com/travel-api/build/internal/config"
	"github.com/travel-api/build/internal/db"
	"github.com/travel-api/build/internal/intent"
	"github.com/travel-api/build/internal/middleware"
	"github.com/travel-api/build/internal/search"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config:%v", err)
	}

	// need to setup database here but will do later when we have some data to persist
	dbCli, err := db.InitDB(cfg.DatabaseURL)
	// Get the underlying sql.DB to manage connection pool and close it on shutdown
	sqlDB, sqlDberr := dbCli.DB()
	if sqlDberr != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", sqlDberr)
	}
	// Ensure the database connection is closed when the application exits
	defer sqlDB.Close()

	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Initialize services and handlers
	searchSvc := search.NewService()
	searchHandler := search.NewHandler(searchSvc)

	intentSvc := intent.NewService(cfg.GeminiAPIKey)
	intentHandler := intent.NewHandler(intentSvc,searchSvc)

	

	approvalSvc := approval.NewService(dbCli, cfg.SendGridAPIKey, cfg.SendGridFrom, cfg.ApproverEmail)
	approvalHandler := approval.NewHandler(approvalSvc)

	// Setup HTTP server and routes here, e.g., using net/http or a router like gorilla/mux

	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.JSONContentTypeMiddleware)

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		middleware.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/intent", intentHandler.HandleIntent).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/search", searchHandler.HandleSearch).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/approval", approvalHandler.HandleApprovalRequest).Methods("POST", "OPTIONS")

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server is running on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully...")
}
