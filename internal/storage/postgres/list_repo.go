package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend_project/internal/domain"
)

var ErrNotFound = errors.New("not found")

// общий таймаут для всех запросов к БД
const dbTimeout = 5 * time.Second

type ListRepo struct {
	pool *pgxpool.Pool
}

func NewListRepo(pool *pgxpool.Pool) *ListRepo {
	return &ListRepo{pool: pool}
}

// Create создаёт новый список
func (r *ListRepo) Create(ctx context.Context, title string) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	id := uuid.New()

	const query = `
		INSERT INTO lists (id, title)
		VALUES ($1, $2)
		RETURNING id, title, created_at
	`

	var list domain.List
	if err := r.pool.QueryRow(ctx, query, id, title).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	); err != nil {
		return domain.List{}, fmt.Errorf("create list: %w", err)
	}

	return list, nil
}

// GetByID получает список по ID
func (r *ListRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
		SELECT id, title, created_at
		FROM lists
		WHERE id = $1
	`

	var list domain.List
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.List{}, ErrNotFound
		}
		return domain.List{}, fmt.Errorf("get list by id: %w", err)
	}

	return list, nil
}

// UpdateTitle обновляет название списка
func (r *ListRepo) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
		UPDATE lists
		SET title = $2
		WHERE id = $1
		RETURNING id, title, created_at
	`

	var list domain.List
	if err := r.pool.QueryRow(ctx, query, id, title).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.List{}, ErrNotFound
		}
		return domain.List{}, fmt.Errorf("update list title: %w", err)
	}

	return list, nil
}

// Delete удаляет список
func (r *ListRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `DELETE FROM lists WHERE id = $1`

	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete list: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// List возвращает списки с пагинацией и их общее количество
func (r *ListRepo) List(ctx context.Context, limit, offset int) ([]domain.List, int, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Общее количество
	const countQuery = `SELECT COUNT(*) FROM lists`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count lists: %w", err)
	}

	// Сами записи
	const query = `
		SELECT id, title, created_at
		FROM lists
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list lists: %w", err)
	}
	defer rows.Close()

	lists := make([]domain.List, 0, limit)
	for rows.Next() {
		var list domain.List
		if err := rows.Scan(&list.ID, &list.Title, &list.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return lists, total, nil
}

// Save делает upsert по id: если запись есть — обновит title, если нет — создаст.
// ВАЖНО: по контракту интерфейса возвращает только error.
func (r *ListRepo) Save(ctx context.Context, l domain.List) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
		INSERT INTO lists (id, title)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET title = EXCLUDED.title
	`

	// Если ID нулевой — сгенерируем новый, чтобы insert был валиден.
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}

	if _, err := r.pool.Exec(ctx, query, l.ID, l.Title); err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}
