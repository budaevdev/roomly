package booking

import "time"

type Listing struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	OwnerID     int64     `json:"owner_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
