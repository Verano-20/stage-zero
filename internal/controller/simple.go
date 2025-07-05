package controller

import (
	"net/http"
	"strconv"

	"github.com/Verano-20/go-crud/internal/model"
	"github.com/Verano-20/go-crud/internal/repository"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type SimpleController struct {
	SimpleRepository *repository.SimpleRepository
}

func NewSimpleController(db *gorm.DB) *SimpleController {
	return &SimpleController{SimpleRepository: repository.NewSimpleRepository(db)}
}

// Create godoc
// @Summary Create a new Simple
// @Description Create a new Simple
// @Tags simple
// @Accept json
// @Produce json
// @Param body body model.SimpleForm true "Simple object"
// @Success 201 {object} model.Simple
// @Router /simple [post]
func (h *SimpleController) Create(c *gin.Context) {
	var simpleForm model.SimpleForm
	if err := c.ShouldBindJSON(&simpleForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	simple, err := h.SimpleRepository.Create(simpleForm.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, simple)
}

// GetAll godoc
// @Summary Get all Simples
// @Description Get all Simples
// @Tags simple
// @Produce json
// @Success 200 {array} model.SimpleDTO
// @Router /simple [get]
func (h *SimpleController) GetAll(c *gin.Context) {
	var simples model.Simples
	simples, err := h.SimpleRepository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, simples.ToDTOs())
}

// GetByID godoc
// @Summary Get Simple by ID
// @Description Get Simple by ID
// @Tags simple
// @Param id path int true "Simple ID"
// @Produce json
// @Success 200 {object} model.SimpleDTO
// @Router /simple/{id} [get]
func (h *SimpleController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	simple, err := h.SimpleRepository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, simple.ToDTO())
}

// Update godoc
// @Summary Update Simple
// @Description Update Simple
// @Tags simple
// @Param id path int true "Simple ID"
// @Accept json
// @Produce json
// @Param body body model.SimpleForm true "Simple object"
// @Success 200 {object} model.SimpleDTO
// @Router /simple/{id} [put]
func (h *SimpleController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var simpleForm model.SimpleForm
	if err := c.ShouldBindJSON(&simpleForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingSimple, err := h.SimpleRepository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	existingSimple.Name = simpleForm.Name

	simple, err := h.SimpleRepository.Update(existingSimple)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, simple.ToDTO())
}

// Delete godoc
// @Summary Delete Simple
// @Description Delete Simple
// @Tags simple
// @Param id path int true "Simple ID"
// @Produce json
// @Success 204
// @Router /simple/{id} [delete]
func (h *SimpleController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	_, err = h.SimpleRepository.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	err = h.SimpleRepository.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
