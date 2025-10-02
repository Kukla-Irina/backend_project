package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend_project/internal/domain"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrValidationTitle = errors.New("title must be 1..100 chars")
)

// интерфейс репозитория
type ListRepo interface {
	Save(ctx context.Context, l domain.List) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.List, error)
	List(ctx context.Context, limit, offset int) ([]domain.List, int, error)
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) (domain.List, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// интерфейс сервиса
type ListService interface {
	Create(ctx context.Context, title string) (domain.List, error)
	Get(ctx context.Context, id uuid.UUID) (domain.List, error)
	List(ctx context.Context, limit, offset int) ([]domain.List, int, error)
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) (domain.List, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// реализация сервиса
type listService struct {
	repo ListRepo
}

func New(repo ListRepo) ListService {
	return &listService{repo: repo}
}

// create
func (s *listService) Create(ctx context.Context, title string) (domain.List, error) {
	if len(title) < domain.TitleMin || len(title) > domain.TitleMax {
		return domain.List{}, ErrValidationTitle
	}
	l := domain.List{
		ID:        uuid.New(),
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Save(ctx, l); err != nil {
		return domain.List{}, err
	}
	return l, nil
}

// get
func (s *listService) Get(ctx context.Context, id uuid.UUID) (domain.List, error) {
	return s.repo.GetByID(ctx, id)
}

// get all
func (s *listService) List(ctx context.Context, limit, offset int) ([]domain.List, int, error) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

// update
func (s *listService) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (domain.List, error) {
	if len(title) < domain.TitleMin || len(title) > domain.TitleMax {
		return domain.List{}, ErrValidationTitle
	}
	return s.repo.UpdateTitle(ctx, id, title)
}

// delete
func (s *listService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
