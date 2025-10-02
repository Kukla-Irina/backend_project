package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"backend_project/internal/http/handlers"
	"backend_project/internal/http/middleware"
	"backend_project/internal/service"
)

// создаем маршруты (эндпоинты)
func NewRouter(svc service.ListService) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	h := handlers.NewListHandlers(svc)

	api.HandleFunc("/lists", h.Create).Methods(http.MethodPost)
	api.HandleFunc("/lists", h.List).Methods(http.MethodGet)
	api.HandleFunc("/lists/{id}", h.Get).Methods(http.MethodGet)
	api.HandleFunc("/lists/{id}", h.Patch).Methods(http.MethodPatch)
	api.HandleFunc("/lists/{id}", h.Delete).Methods(http.MethodDelete)

	return r
}
