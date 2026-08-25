package booking

import "time"

type Booking struct {
	ID        int64     `json:"id"`
	ListingID int64     `json:"listing_id"`
	GuestID   int64     `json:"guest_id"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
