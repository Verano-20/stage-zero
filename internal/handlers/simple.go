package handlers

import (
	"net/http"

	"github.com/Verano-20/go-crud/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SimpleHandler struct {
	DB *gorm.DB
}

func NewSimpleHandler(db *gorm.DB) *SimpleHandler {
	return &SimpleHandler{DB: db}
}

// Create godoc
// @Summary Create a new Simple
// @Description Create a new Simple
// @Tags simple
// @Accept json
// @Produce json
// @Param body body models.SimpleForm true "Simple object"
// @Success 201 {object} models.Simple
// @Router /simple [post]
func (h *SimpleHandler) Create(c *gin.Context) {
	var simple models.Simple
	if err := c.ShouldBindJSON(&simple); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&simple).Error; err != nil {
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
// @Success 200 {array} models.Simple
// @Router /simple [get]
func (h *SimpleHandler) GetAll(c *gin.Context) {
	var simples []models.Simple
	if err := h.DB.Find(&simples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, simples)
}

// GetByID godoc
// @Summary Get Simple by ID
// @Description Get Simple by ID
// @Tags simple
// @Param id path int true "Simple ID"
// @Produce json
// @Success 200 {object} models.Simple
// @Router /simple/{id} [get]
func (h *SimpleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var simple models.Simple

	if err := h.DB.First(&simple, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusFound, simple)
}

// Update godoc
// @Summary Update Simple
// @Description Update Simple
// @Tags simple
// @Param id path int true "Simple ID"
// @Accept json
// @Produce json
// @Param body body models.SimpleForm true "Simple object"
// @Success 200 {object} models.Simple
// @Router /simple/{id} [put]
func (h *SimpleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var simple models.Simple

	if err := h.DB.First(&simple, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := c.ShouldBindJSON(&simple); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&simple).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, simple)
}

// Delete godoc
// @Summary Delete Simple
// @Description Delete Simple
// @Tags simple
// @Param id path int true "Simple ID"
// @Produce json
// @Success 204
// @Router /simple/{id} [delete]
func (h *SimpleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var simple models.Simple

	if err := h.DB.First(&simple, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := h.DB.Delete(&simple).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
