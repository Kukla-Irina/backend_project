package storage

import "backend_project/internal/domain"

// ListRepository — интерфейс для работы со списками
type ListRepository interface {
	Create(title string) (domain.List, error)
	GetByID(id string) (domain.List, error)
	UpdateTitle(id, title string) (domain.List, error)
	Delete(id string) error
	List(limit, offset int) ([]domain.List, int, error)
}
