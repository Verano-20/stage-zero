package controller

import (
	"net/http"
	"strconv"

	"github.com/Verano-20/stage-zero/internal/logger"
	"github.com/Verano-20/stage-zero/internal/model"
	"github.com/Verano-20/stage-zero/internal/response"
	"github.com/Verano-20/stage-zero/internal/service"
	"github.com/Verano-20/stage-zero/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SimpleController struct {
	SimpleService service.SimpleService
}

func NewSimpleController(simpleService service.SimpleService) *SimpleController {
	return &SimpleController{SimpleService: simpleService}
}

// Create godoc
// @Summary Create a new Simple
// @Description Create a new Simple with the provided details. The name field is required and must be a non-empty string.
// @Tags Simple
// @Accept json
// @Produce json
// @Param simple body model.SimpleForm true "Simple details"
// @Success 201 {object} response.ApiResponse "Simple created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request format"
// @Failure 500 {object} response.ErrorResponse "Internal server error during resource creation"
// @Router /simple [post]
func (c *SimpleController) Create(ctx *gin.Context) {
	var simpleForm model.SimpleForm
	if err := ctx.ShouldBindJSON(&simpleForm); err != nil {
		utils.HandleBindingErrors(ctx, err, "create")
		return
	}

	simple, err := c.SimpleService.CreateSimple(ctx, simpleForm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to create Simple"})
		return
	}

	ctx.JSON(http.StatusCreated, response.ApiResponse{Message: "Simple created successfully", Data: simple.ToDTO()})
}

// GetAll godoc
// @Summary Get all Simples
// @Description Get all Simples. Returns an array of Simple objects. Returns an empty array if none exist.
// @Tags Simple
// @Produce json
// @Success 200 {object} response.ApiResponse "Simples retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving Simples"
// @Router /simple [get]
func (c *SimpleController) GetAll(ctx *gin.Context) {
	simples, err := c.SimpleService.GetAllSimples(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to retrieve Simples"})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simples retrieved successfully", Data: simples.ToDTOs()})
}

// GetByID godoc
// @Summary Get Simple by ID
// @Description Find a Simple by its unique ID
// @Tags Simple
// @Param id path int true "Simple ID"
// @Produce json
// @Success 200 {object} response.ApiResponse "Simple retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid ID format or value"
// @Failure 404 {object} response.ErrorResponse "Simple not found"
// @Router /simple/{id} [get]
func (c *SimpleController) GetByID(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for get by id", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	simple, err := c.SimpleService.GetSimpleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple retrieved successfully", Data: simple.ToDTO()})
}

// Update godoc
// @Summary Update an existing Simple
// @Description Update a Simple identified by its ID with new data. The ID must exist and the request body must contain valid data.
// @Tags Simple
// @Accept json
// @Produce json
// @Param id path int true "Simple ID to update"
// @Param simple body model.SimpleForm true "Updated Simple details"
// @Success 200 {object} response.ApiResponse "Simple updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request body format"
// @Failure 404 {object} response.ErrorResponse "Simple not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error during update operation"
// @Router /simple/{id} [put]
func (c *SimpleController) Update(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for update", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	var simpleForm model.SimpleForm
	if err := ctx.ShouldBindJSON(&simpleForm); err != nil {
		utils.HandleBindingErrors(ctx, err, "update")
		return
	}

	existingSimple, err := c.SimpleService.GetSimpleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	simple, err := c.SimpleService.UpdateSimple(ctx, existingSimple, simpleForm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to update Simple"})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple updated successfully", Data: simple.ToDTO()})
}

// Delete godoc
// @Summary Delete a Simple
// @Description Permanently delete a Simple identified by its ID. This operation cannot be undone.
// @Tags Simple
// @Produce json
// @Param id path int true "Simple ID to delete"
// @Success 200 {object} response.ApiResponse "Simple deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid ID format or value"
// @Failure 404 {object} response.ErrorResponse "Simple not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error during deletion"
// @Router /simple/{id} [delete]
func (c *SimpleController) Delete(ctx *gin.Context) {
	log := logger.GetFromContext(ctx)

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warn("Invalid ID format for delete", zap.String("id_param", idParam), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	existingSimple, err := c.SimpleService.GetSimpleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Simple not found"})
		return
	}

	err = c.SimpleService.DeleteSimple(ctx, existingSimple)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to delete Simple"})
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{Message: "Simple deleted successfully", Data: nil})
}
