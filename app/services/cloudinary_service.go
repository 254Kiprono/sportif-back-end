package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"webuye-sportif/app/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// ImageFolder defines the folder/context the image belongs to.
type ImageFolder string

const (
	FolderNews         ImageFolder = "sportif/news"
	FolderMatchPreview ImageFolder = "sportif/fixtures/previews"
	FolderMatchPhotos  ImageFolder = "sportif/fixtures/match_photos"
	FolderGeneral      ImageFolder = "sportif/general"
)

// UploadResult holds the returned data after a successful upload.
type UploadResult struct {
	PublicID     string `json:"public_id"`
	SecureURL    string `json:"secure_url"`
	Format       string `json:"format"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Bytes        int    `json:"bytes"`
	OriginalName string `json:"original_name"`
}

type CloudinaryService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader, folder ImageFolder) (*UploadResult, error)
	DeleteImage(publicID string) error
}

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cfg *config.Config) (CloudinaryService, error) {
	if cfg.CloudinaryName == "" || cfg.CloudinaryKey == "" || cfg.CloudinarySecret == "" {
		return nil, fmt.Errorf("cloudinary credentials are not configured in environment variables")
	}

	cld, err := cloudinary.NewFromParams(cfg.CloudinaryName, cfg.CloudinaryKey, cfg.CloudinarySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &cloudinaryService{cld: cld}, nil
}

// UploadImage uploads a file to Cloudinary with smart compression.
//
// Transformation breakdown (no quality loss strategy):
//   - c_limit,w_2000,h_2000 → Constrain max dimensions only; never upscales.
//   - q_auto:best            → Cloudinary picks the optimal quality level.
//     "auto:best" prioritizes visual quality — typically ~40-60% file size
//     reduction with NO perceptible quality difference to the human eye.
//   - f_auto                 → Serves WebP to modern browsers, JPEG/PNG to older
//     ones — all automatically based on browser Accept headers.
//
// This is Cloudinary's recommended approach for "high quality, small size".
func (s *cloudinaryService) UploadImage(
	file multipart.File,
	header *multipart.FileHeader,
	folder ImageFolder,
) (*UploadResult, error) {
	// Validate file is an image
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("invalid file type '%s': only jpg, jpeg, png, webp and gif are allowed", ext)
	}

	// Enforce max file size: 10MB
	const maxFileSize = 10 * 1024 * 1024
	if header.Size > maxFileSize {
		return nil, fmt.Errorf("file size %d bytes exceeds the maximum of 10MB", header.Size)
	}

	// Generate a unique public ID using timestamp
	timestamp := time.Now().UnixMilli()
	cleanName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	publicID := fmt.Sprintf("%s/%s_%d", folder, cleanName, timestamp)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Transformation string: constrain size + auto best quality + auto format (WebP/JPEG/PNG)
	transformation := "c_limit,w_2000,h_2000/q_auto:best/f_auto"

	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicID,
		Transformation: transformation,
		Tags:           api.CldAPIArray{string(folder), "sportif-fc"},
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return &UploadResult{
		PublicID:     resp.PublicID,
		SecureURL:    resp.SecureURL,
		Format:       resp.Format,
		Width:        resp.Width,
		Height:       resp.Height,
		Bytes:        resp.Bytes,
		OriginalName: header.Filename,
	}, nil
}

// DeleteImage removes an image from Cloudinary by its public ID.
func (s *cloudinaryService) DeleteImage(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
