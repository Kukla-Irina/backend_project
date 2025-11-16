// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	apihttp "backend_project/internal/http"
// 	"backend_project/internal/service"
// 	"backend_project/internal/storage/mem"
// )

// // запуск сервиса

// func main() {
// 	repo := mem.NewRepo()
// 	svc := service.New(repo)
// 	router := apihttp.NewRouter(svc)

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	addr := ":" + port

// 	fmt.Println("Server listening on http://localhost" + addr)
// 	log.Fatal(http.ListenAndServe(addr, router))
// }

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend_project/internal/config"
	"backend_project/internal/database"
	httphandlers "backend_project/internal/http"
	"backend_project/internal/service"
	"backend_project/internal/storage/postgres"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем контекст для работы
	ctx := context.Background()

	// Подключаемся к PostgreSQL
	log.Println("Connecting to database...")
	pool, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	// Создаем репозиторий PostgreSQL
	repo := postgres.NewListRepo(pool)

	// Создаем сервис
	svc := service.New(repo)

	// Создаем HTTP-роутер
	router := httphandlers.NewRouter(svc)

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
