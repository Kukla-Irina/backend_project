package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	apihttp "backend_project/internal/http"
	"backend_project/internal/service"
	"backend_project/internal/storage/mem"
)

// запуск сервиса

func main() {
	repo := mem.NewRepo()
	svc := service.New(repo)
	router := apihttp.NewRouter(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Println("Server listening on http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
