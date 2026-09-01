package api

import (
	"errors"
	"net/http"

	"github.com/budaevdev/roomly/internal/auth"
	"github.com/budaevdev/roomly/internal/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	u := storage.User{Username: req.Username, PasswordHash: hash}
	err = storage.CreateUser(s.db, &u)
	if err != nil {
		if errors.Is(err, storage.ErrUsernameTaken) {
			c.String(http.StatusConflict, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	u, err := storage.GetUserByUsername(s.db, req.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			c.String(http.StatusUnauthorized, "invalid username or password")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.String(http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := auth.CreateSession(c.Request.Context(), s.rdb, u.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("session", token, 0, "/", "", false, true)
	c.Status(http.StatusOK)
}

func (s *Server) Logout(c *gin.Context) {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	auth.DeleteSession(c.Request.Context(), s.rdb, cookie)

	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Status(http.StatusNoContent)
}
