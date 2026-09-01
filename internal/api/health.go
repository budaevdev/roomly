package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
