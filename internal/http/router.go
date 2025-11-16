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

	// health — подгоняем под OpenAPI (JSON {"status":"ok"})
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Swagger: спецификация (yaml)
	r.HandleFunc("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		http.ServeFile(w, r, "docs/openapi.yaml") // файл с твоей OpenAPI-спекой
	}).Methods(http.MethodGet)

	// Swagger: UI
	r.HandleFunc("/swagger", swaggerUIHandler).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	h := handlers.NewListHandlers(svc)

	api.HandleFunc("/lists", h.Create).Methods(http.MethodPost)
	api.HandleFunc("/lists", h.List).Methods(http.MethodGet)
	api.HandleFunc("/lists/{id}", h.Get).Methods(http.MethodGet)
	api.HandleFunc("/lists/{id}", h.Patch).Methods(http.MethodPatch)
	api.HandleFunc("/lists/{id}", h.Delete).Methods(http.MethodDelete)

	return r
}

// хендлер, который отдаёт страницу со Swagger UI
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerHTML))
}

// простая HTML-страница, которая поднимает Swagger UI из CDN
const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Lists API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({
        url: '/swagger/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`
