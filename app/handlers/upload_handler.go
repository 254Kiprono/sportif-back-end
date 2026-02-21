package handlers

import (
	"net/http"

	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles all image uploads via Cloudinary.
// It is separate from specific domain handlers (news, fixtures) so it's reusable.
type UploadHandler struct {
	cloudinaryService services.CloudinaryService
}

func NewUploadHandler(cloudinaryService services.CloudinaryService) *UploadHandler {
	return &UploadHandler{cloudinaryService: cloudinaryService}
}

// UploadNewsImage handles image upload for news articles.
// Protected: Author + Admin only (they create news content).
//
// Request: multipart/form-data with field "image"
// Response: { "url": "...", "public_id": "...", "width": ..., "height": ..., "bytes": ... }
func (h *UploadHandler) UploadNewsImage(c *gin.Context) {
	h.handleUpload(c, services.FolderNews)
}

// UploadMatchPreview handles pre-match preview photo uploads.
// Protected: CX + Admin only.
//
// Request: multipart/form-data with field "image"
func (h *UploadHandler) UploadMatchPreview(c *gin.Context) {
	h.handleUpload(c, services.FolderMatchPreview)
}

// UploadMatchPhoto handles post-match action photo uploads.
// Protected: CX + Admin only.
//
// Request: multipart/form-data with field "image"
func (h *UploadHandler) UploadMatchPhoto(c *gin.Context) {
	h.handleUpload(c, services.FolderMatchPhotos)
}

// handleUpload is the shared core upload logic.
func (h *UploadHandler) handleUpload(c *gin.Context, folder services.ImageFolder) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No image file provided. Use form-data field 'image'.",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer file.Close()

	result, err := h.cloudinaryService.UploadImage(file, fileHeader, folder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":           result.SecureURL,
		"public_id":     result.PublicID,
		"width":         result.Width,
		"height":        result.Height,
		"bytes":         result.Bytes,
		"format":        result.Format,
		"original_name": result.OriginalName,
	})
}

// DeleteImage removes an image from Cloudinary.
// Protected: Admin only.
//
// Request JSON: { "public_id": "sportif/news/my-image_1234567890" }
func (h *UploadHandler) DeleteImage(c *gin.Context) {
	var input struct {
		PublicID string `json:"public_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cloudinaryService.DeleteImage(input.PublicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully", "public_id": input.PublicID})
}
