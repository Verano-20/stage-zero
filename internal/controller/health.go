package controller

import (
	"net/http"

	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/Verano-20/stage-zero/internal/response"
	"github.com/gin-gonic/gin"
)

// GetHealth godoc
// @Summary Get health
// @Description Get health of the server
// @Tags Health
// @Produce json
// @Success 200 {object} response.ApiResponse "Server is healthy and operational" example({"message": "OK", "data": null})
// @Router /health [get]
func GetHealth(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	log.Debug("Health check requested...")

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "OK"})

	log.Debug("Health check completed")
}
