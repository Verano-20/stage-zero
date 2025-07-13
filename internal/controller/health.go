package controller

import (
	"net/http"

	"github.com/Verano-20/go-crud/internal/response"
	"github.com/gin-gonic/gin"
)

// GetHealth godoc
// @Summary Get health
// @Description Get health of the server
// @Tags health
// @Produce json
// @Success 200 {object} response.ApiResponse "Server is healthy and operational" example({"message": "OK", "data": null})
// @Router /health [get]
func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, response.ApiResponse{Message: "OK"})
}
