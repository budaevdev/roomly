package storage

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/budaevdev/roomly/internal/booking"
	"github.com/joho/godotenv"
)

func TestCreateBooking_ConcurrentOverlap(t *testing.T) {
	godotenv.Load("../../.env")

	dbURL := os.Getenv("DB_URL")

	conn, err := Connect(dbURL)

	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	var guestID int64
	err = conn.QueryRow(`INSERT INTO users (username) VALUES ($1) RETURNING id`, fmt.Sprintf("test-guest-%d", time.Now().UnixNano())).Scan(&guestID)

	if err != nil {
		t.Fatal(err)
	}

	l := booking.Listing{Title: "Test Listing", Description: "for overlap test", OwnerID: guestID}

	err = CreateListing(conn, &l)

	if err != nil {
		t.Fatal(err)
	}

	const n = 10

	var wg sync.WaitGroup
	errs := make([]error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			b := &booking.Booking{
				ListingID: l.ID,
				GuestID:   guestID,
				StartDate: "2026-10-01",
				EndDate:   "2026-10-05",
			}
			errs[i] = CreateBooking(conn, b)
		}(i)
	}
	wg.Wait()

	var successes, rejected int
	for _, err := range errs {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrOverlappingBooking) {
			rejected++
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successes != 1 || rejected != n-1 {
		t.Fatalf("successes: %d, rejected: %d", successes, rejected)
	}
}
