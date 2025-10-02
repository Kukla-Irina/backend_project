package domain

import (
	"time"

	"github.com/google/uuid"
)

// константы для максимального и минимального количества символов
const (
	TitleMin = 1
	TitleMax = 100
)

// структура списка
type List struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
