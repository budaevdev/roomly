package api

import (
	"database/sql"

	"github.com/budaevdev/roomly/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewServer(db *sql.DB, rdb *redis.Client) *Server {
	return &Server{db: db, rdb: rdb}
}

func (s *Server) Routes() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", s.Healthz)

	r.POST("/listings", auth.RequireAuth(s.rdb), s.CreateListing)
	r.GET("/listings", s.GetListings)
	r.GET("/listings/search", s.SearchListings)

	r.POST("/bookings", auth.RequireAuth(s.rdb), s.CreateBooking)
	r.PATCH("/bookings/:id", s.CancelBooking)

	r.POST("/register", s.Register)
	r.POST("/login", s.Login)
	r.POST("/logout", s.Logout)

	return r
}
