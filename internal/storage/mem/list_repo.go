package mem

// репозиторий

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"

	"backend_project/internal/domain"
	"backend_project/internal/service"
)

type Repo struct {
	mu    sync.RWMutex
	data  map[uuid.UUID]domain.List
	order []uuid.UUID
}

func NewRepo() *Repo {
	return &Repo{
		data:  make(map[uuid.UUID]domain.List),
		order: make([]uuid.UUID, 0),
	}
}

func (r *Repo) Save(ctx context.Context, l domain.List) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[l.ID]; !ok {
		r.order = append(r.order, l.ID)
	}
	r.data[l.ID] = l
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (domain.List, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.data[id]
	if !ok {
		return domain.List{}, service.ErrNotFound
	}
	return l, nil
}

func (r *Repo) List(ctx context.Context, limit, offset int) ([]domain.List, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.order)
	if offset > total {
		return []domain.List{}, total, nil
	}

	items := make([]domain.List, 0, total)
	for _, id := range r.order {
		items = append(items, r.data[id])
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	// Пагинация
	end := offset + limit
	if limit == 0 || end > total {
		end = total
	}
	if offset < 0 {
		offset = 0
	}
	return items[offset:end], total, nil
}

func (r *Repo) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (domain.List, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	l, ok := r.data[id]
	if !ok {
		return domain.List{}, service.ErrNotFound
	}
	l.Title = title
	r.data[id] = l
	return l, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return service.ErrNotFound
	}
	delete(r.data, id)
	for i, v := range r.order {
		if v == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}
