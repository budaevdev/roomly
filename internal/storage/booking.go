package storage

import (
	"database/sql"
	"errors"

	"github.com/budaevdev/roomly/internal/booking"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrOverlappingBooking = errors.New("listing not available for these dates")
var ErrBookingNotFound = errors.New("booking not found or already cancelled")

func CreateBooking(db *sql.DB, b *booking.Booking) error {
	const query = `INSERT INTO bookings (listing_id, guest_id, during, status) VALUES ($1, $2, daterange($3, $4, '[)'), 'active') RETURNING id, created_at, updated_at, status`

	err := db.QueryRow(query, b.ListingID, b.GuestID, b.StartDate, b.EndDate).Scan(&b.ID, &b.CreatedAt,
		&b.UpdatedAt, &b.Status)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			return ErrOverlappingBooking
		}
		return err
	}

	return nil
}

func CancelBooking(db *sql.DB, id int64) error {
	const query = `UPDATE bookings SET status = 'cancelled', updated_at = now() WHERE id = $1 AND status = 'active'`

	res, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookingNotFound
	}

	return nil
}
