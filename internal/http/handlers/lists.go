package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"backend_project/internal/service"
)

// унифицированная ошибка
type apiError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// для одинакового вывода JSON ответа
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// для разделения слоев
type ListHandlers struct {
	svc service.ListService
}

// собираем набор хэндлеров с подключённым сервисом
func NewListHandlers(svc service.ListService) *ListHandlers {
	return &ListHandlers{svc: svc}
}

// модель тела запроса на создание списка
type createListReq struct {
	Title string `json:"title"`
}

// модель тела запроса на обновление
type updateListReq struct {
	Title *string `json:"title"`
}

// POST /api/v1/lists
func (h *ListHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req createListReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "invalid json", nil})
		return
	}
	l, err := h.svc.Create(r.Context(), req.Title)
	if err != nil {
		if err == service.ErrValidationTitle {
			writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", err.Error(), nil})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiError{"INTERNAL_ERROR", "internal error", nil})
		return
	}
	writeJSON(w, http.StatusCreated, l)
}

// GET /api/v1/lists?limit=&offset=
func (h *ListHandlers) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit == 0 && q.Get("limit") == "" {
		limit = 20
	}
	items, total, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{"INTERNAL_ERROR", "internal error", nil})
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	writeJSON(w, http.StatusOK, items)
}

// GET /api/v1/lists/{id}
func (h *ListHandlers) Get(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "invalid id", nil})
		return
	}
	l, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			writeJSON(w, http.StatusNotFound, apiError{"NOT_FOUND", "not found", nil})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiError{"INTERNAL_ERROR", "internal error", nil})
		return
	}
	writeJSON(w, http.StatusOK, l)
}

// PATCH /api/v1/lists/{id}
func (h *ListHandlers) Patch(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "invalid id", nil})
		return
	}
	var req updateListReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "invalid json", nil})
		return
	}
	if req.Title == nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "title is required", nil})
		return
	}
	l, err := h.svc.UpdateTitle(r.Context(), id, *req.Title)
	if err != nil {
		switch err {
		case service.ErrValidationTitle:
			writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", err.Error(), nil})
		case service.ErrNotFound:
			writeJSON(w, http.StatusNotFound, apiError{"NOT_FOUND", "not found", nil})
		default:
			writeJSON(w, http.StatusInternalServerError, apiError{"INTERNAL_ERROR", "internal error", nil})
		}
		return
	}
	writeJSON(w, http.StatusOK, l)
}

// DELETE /api/v1/lists/{id}
func (h *ListHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{"VALIDATION_FAILED", "invalid id", nil})
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if err == service.ErrNotFound {
			writeJSON(w, http.StatusNotFound, apiError{"NOT_FOUND", "not found", nil})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiError{"INTERNAL_ERROR", "internal error", nil})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
