package storage

import (
	"database/sql"

	"github.com/budaevdev/roomly/internal/booking"
)

func CreateListing(db *sql.DB, l *booking.Listing) error {
	const query = `INSERT INTO listings (title, description, owner_id) VALUES ($1, $2, $3) 
  RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, l.Title, l.Description, l.OwnerID).Scan(&l.ID, &l.CreatedAt,
		&l.UpdatedAt)
	return err
}

func GetListings(db *sql.DB) ([]booking.Listing, error) {
	const query = `SELECT id, title, description, owner_id, created_at, updated_at FROM listings`
	listings := []booking.Listing{}

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l booking.Listing
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.OwnerID, &l.CreatedAt,
			&l.UpdatedAt)
		if err != nil {
			return nil, err
		}

		listings = append(listings, l)
	}
	return listings, rows.Err()
}

func GetAvailableListings(db *sql.DB, start, end string) ([]booking.Listing, error) {
	const query = `
                SELECT l.id, l.title, l.description, l.owner_id, l.created_at, l.updated_at
                FROM listings l
                WHERE NOT EXISTS (
                        SELECT 1 FROM bookings b
                        WHERE b.listing_id = l.id
                          AND b.status = 'active'
                          AND b.during && daterange($1, $2, '[)')
                )`

	listings := []booking.Listing{}

	rows, err := db.Query(query, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l booking.Listing
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.OwnerID, &l.CreatedAt,
			&l.UpdatedAt)
		if err != nil {
			return nil, err
		}

		listings = append(listings, l)
	}
	return listings, rows.Err()
}
