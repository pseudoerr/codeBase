package models

import "time"

type Mission struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"` // Внешний ключ на пользователя
	Title       string    `json:"title"`
	Points      int       `json:"points"` // Количество очков за выполнение задания
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
