package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"webuye-sportif/app/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ImageFolder defines the folder/context the image belongs to.
type ImageFolder string

const (
	FolderNews         ImageFolder = "sportif/news"
	FolderMatchPreview ImageFolder = "sportif/fixtures/previews"
	FolderMatchPhotos  ImageFolder = "sportif/fixtures/match_photos"
	FolderGeneral      ImageFolder = "sportif/general"
	FolderJerseys      ImageFolder = "sportif/store/jerseys"
	FolderPlayers      ImageFolder = "sportif/players"
)

// UploadResult holds the returned data after a successful upload.
type UploadResult struct {
	PublicID     string `json:"public_id"`
	SecureURL    string `json:"secure_url"`
	Format       string `json:"format"`
	OriginalName string `json:"original_name"`
	Bytes        int    `json:"bytes"`
}

type StorageService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader, folder ImageFolder) (*UploadResult, error)
	DeleteImage(publicID string) error
}

type b2StorageService struct {
	client     *s3.Client
	bucketName string
	endpoint   string
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	if cfg.B2KeyID == "" || cfg.B2ApplicationKey == "" || cfg.B2BucketName == "" || cfg.B2Endpoint == "" || cfg.B2Region == "" {
		return nil, fmt.Errorf("Backblaze B2 credentials are not fully configured in environment variables")
	}

	keyID := strings.TrimSpace(cfg.B2KeyID)
	appKey := strings.TrimSpace(cfg.B2ApplicationKey)
	region := strings.TrimSpace(cfg.B2Region)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(keyID, appKey, "")),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load B2 AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.B2Endpoint)
		o.UsePathStyle = true
		o.Region = region
	})

	return &b2StorageService{
		client:     client,
		bucketName: cfg.B2BucketName,
		endpoint:   cfg.B2Endpoint,
	}, nil
}

func (s *b2StorageService) UploadImage(
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

	// Detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("could not read file to detect content type")
	}
	// Reset file pointer
	file.Seek(0, 0)
	contentType := http.DetectContentType(buffer)

	// Generate a unique public ID using timestamp
	timestamp := time.Now().UnixMilli()
	cleanName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	publicID := fmt.Sprintf("%s/%s_%d%s", folder, cleanName, timestamp, ext)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(publicID),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(header.Size),
		// Use unsigned payload to avoid "Seed signature is invalid" on B2
	})
	if err != nil {
		return nil, fmt.Errorf("B2 upload failed (Key: %s): %w", publicID, err)
	}

	// Construct public URL
	cleanEndpoint := strings.TrimSuffix(s.endpoint, "/")
	secureURL := fmt.Sprintf("%s/%s/%s", cleanEndpoint, s.bucketName, publicID)

	return &UploadResult{
		PublicID:     publicID,
		SecureURL:    secureURL,
		Format:       strings.TrimPrefix(ext, "."),
		OriginalName: header.Filename,
		Bytes:        int(header.Size),
	}, nil
}

// DeleteImage removes an image from B2 by its public ID.
func (s *b2StorageService) DeleteImage(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(publicID),
	})
	return err
}
