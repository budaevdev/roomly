package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/budaevdev/roomly/internal/auth"
	"github.com/budaevdev/roomly/internal/booking"
	"github.com/budaevdev/roomly/internal/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateListing(c *gin.Context) {
	var l booking.Listing

	if err := c.ShouldBindJSON(&l); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ownerID, _ := auth.UserIDFromContext(c.Request.Context())
	l.OwnerID = ownerID

	err := storage.CreateListing(s.db, &l)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, l)
}

func (s *Server) GetListings(c *gin.Context) {
	listings, err := storage.GetListings(s.db)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listings)
}

func (s *Server) SearchListings(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	if start == "" || end == "" {
		c.String(http.StatusBadRequest, "start and end query params are required")
		return
	}

	key := fmt.Sprintf("avail:%s:%s", start, end)

	val, err := s.rdb.Get(c.Request.Context(), key).Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}

	listings, err := storage.GetAvailableListings(s.db, start, end)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := json.Marshal(listings)
	if err == nil {
		s.rdb.Set(c.Request.Context(), key, data, 30*time.Second)
	}

	c.JSON(http.StatusOK, listings)
}
