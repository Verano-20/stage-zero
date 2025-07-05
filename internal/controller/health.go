package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealth godoc
// @Summary Get health
// @Description Get health of the server
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
