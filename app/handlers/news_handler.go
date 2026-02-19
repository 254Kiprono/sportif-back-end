package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	service services.NewsService
}

func NewNewsHandler(service services.NewsService) *NewsHandler {
	return &NewsHandler{service}
}

func (h *NewsHandler) GetAll(c *gin.Context) {
	// Public can only see published news
	news, err := h.service.GetNews(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) GetAllAdmin(c *gin.Context) {
	// Admin can see all news
	news, err := h.service.GetNews(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	news, err := h.service.GetNewsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) Create(c *gin.Context) {
	authorID, _ := c.Get("user_id")
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateNews(&news, authorID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, news)
}

func (h *NewsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateNews(id, &news); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteNews(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "News deleted"})
}
