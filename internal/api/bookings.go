package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/budaevdev/roomly/internal/auth"
	"github.com/budaevdev/roomly/internal/booking"
	"github.com/budaevdev/roomly/internal/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateBooking(c *gin.Context) {
	var b booking.Booking

	if err := c.ShouldBindJSON(&b); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	guestID, _ := auth.UserIDFromContext(c.Request.Context())
	b.GuestID = guestID

	err := storage.CreateBooking(s.db, &b)
	if err != nil {
		if errors.Is(err, storage.ErrOverlappingBooking) {
			c.String(http.StatusConflict, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, b)
}

func (s *Server) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = storage.CancelBooking(s.db, id)
	if err != nil {
		if errors.Is(err, storage.ErrBookingNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
