package main

import (
	"log"

	"github.com/takagi_hisashi/go-best-practice/web-api/config"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/database"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/database/repository"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/server"
	infraHTTP "github.com/takagi_hisashi/go-best-practice/web-api/internal/infrastructure/http"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/handler"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/router"
	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/gateway/jsonplaceholder"
	postUseCase "github.com/takagi_hisashi/go-best-practice/web-api/internal/usecase/post"
	userUseCase "github.com/takagi_hisashi/go-best-practice/web-api/internal/usecase/user"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed initial data
	if err := database.SeedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	// Setup infrastructure
	httpClient := infraHTTP.NewHTTPClient()

	// Setup repositories (database-based)
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Setup fallback gateways for external API (optional)
	postGateway := jsonplaceholder.NewPostGateway(cfg.JSONPlaceholderURL, httpClient)
	userGateway := jsonplaceholder.NewUserGateway(cfg.JSONPlaceholderURL, httpClient)
	_ = postGateway // Use if needed for external data
	_ = userGateway // Use if needed for external data

	// Setup use cases with database repositories
	postService := postUseCase.NewService(postRepo)
	userService := userUseCase.NewService(userRepo)

	// Setup handlers
	postHandler := handler.NewPostHandler(postService)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := router.NewRouter(postHandler, userHandler)
	mux := router.Setup()

	// Start server
	srv := server.NewServer(cfg.ServerPort)
	srv.Run(mux)
}