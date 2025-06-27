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

func CreateSimpleHandler(db *gorm.DB) *SimpleHandler {
	return &SimpleHandler{DB: db}
}

// Create a new Simple
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

// Get all Simples
func (h *SimpleHandler) GetAll(c *gin.Context) {
	var simples []models.Simple
	if err := h.DB.Find(&simples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, simples)
}

// Get Simple by ID
func (h *SimpleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var simple models.Simple

	if err := h.DB.First(&simple, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusFound, simple)
}

// Update Simple
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

// Delete Simple
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
