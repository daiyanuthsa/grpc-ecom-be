package entity

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID
	UserID    string
	ProductID string
	Quantity  int
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}
