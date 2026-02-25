package services

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
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
	"github.com/disintegration/imaging"
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
	FolderTickets      ImageFolder = "sportif/tickets"
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
	UploadImage(file multipart.File, header *multipart.FileHeader, folder ImageFolder, customPrefix string) (*UploadResult, error)
	UploadData(data []byte, fileName string, contentType string, folder ImageFolder) (*UploadResult, error)
	DeleteImage(publicID string) error
}

type r2StorageService struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	if cfg.R2Endpoint == "" || cfg.R2AccessKeyID == "" || cfg.R2SecretKey == "" || cfg.R2BucketName == "" {
		return nil, fmt.Errorf("R2_ENDPOINT, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, and R2_BUCKET_NAME must be set")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.R2Endpoint,
			HostnameImmutable: true,
			SigningRegion:     "auto",
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretKey, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &r2StorageService{
		client:     client,
		bucketName: strings.TrimSpace(cfg.R2BucketName),
		publicURL:  strings.TrimSuffix(strings.TrimSpace(cfg.R2PublicURL), "/"),
	}, nil
}

func (s *r2StorageService) UploadImage(file multipart.File, header *multipart.FileHeader, folder ImageFolder, customPrefix string) (*UploadResult, error) {
	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("invalid file type '%s': only jpg, jpeg, png, webp and gif", ext)
	}
	const maxFileSize = 10 * 1024 * 1024
	if header.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds the 10MB limit")
	}

	// Read entire file into memory to detect content type
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read uploaded file: %w", err)
	}

	// Detect content type using http.DetectContentType
	contentType := "application/octet-stream"
	if len(fileBytes) > 512 {
		contentType = http.DetectContentType(fileBytes[:512])
	} else if len(fileBytes) > 0 {
		contentType = http.DetectContentType(fileBytes)
	}

	// Seek back to start before reading for image processing
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("could not seek file: %w", err)
	}

	var uploadBody io.Reader = file
	var uploadSize = header.Size

	// Process image if it is JPEG or PNG
	if contentType == "image/jpeg" || contentType == "image/png" {
		img, err := imaging.Decode(file)
		if err == nil {
			// Resize while preserving aspect ratio if width is larger than 1200px
			if img.Bounds().Dx() > 1200 {
				img = imaging.Resize(img, 1200, 0, imaging.Lanczos)
			}

			// Encode to WEBP wrapper buffer or JPEG with high compression
			var buf bytes.Buffer
			if contentType == "image/jpeg" {
				err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(75))
			} else {
				// Keep PNG intact, or could also compress PNG, but typical use case just resize
				err = imaging.Encode(&buf, img, imaging.PNG, imaging.PNGCompressionLevel(jpeg.DefaultQuality))
			}

			if err == nil {
				// We successfully compressed it
				uploadBody = bytes.NewReader(buf.Bytes())
				uploadSize = int64(buf.Len())
				fileBytes = buf.Bytes() // For returning byte size later
			} else {
				// If compression failed, just fallback to the original file
				if _, err := file.Seek(0, io.SeekStart); err != nil {
					return nil, fmt.Errorf("could not seek file after compression failed: %w", err)
				}
			}
		} else {
			// If decoding fails, fall back to original
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				return nil, fmt.Errorf("could not seek file after decode failure: %w", err)
			}
		}
	} else {
		// Seek back to start for the upload
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("could not seek file: %w", err)
		}
	}

	// Build unique file path
	timestamp := time.Now().UnixMilli()
	var cleanName string
	if customPrefix != "" {
		cleanName = customPrefix
	} else {
		cleanName = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	}

	// Sanitize name: replace spaces and non-alphanumeric chars with underscores
	cleanName = strings.ToLower(cleanName)
	cleanName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, cleanName)

	fileName := fmt.Sprintf("%s/%s_%d%s", folder, cleanName, timestamp, ext)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fileName),
		Body:          uploadBody,
		ContentLength: aws.Int64(uploadSize),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("R2 upload failed: %w", err)
	}

	// Construct public download URL
	var secureURL string
	if s.publicURL != "" {
		secureURL = fmt.Sprintf("%s/%s", s.publicURL, fileName)
	} else {
		baseEndpoint := ""
		if s.client.Options().BaseEndpoint != nil {
			baseEndpoint = *s.client.Options().BaseEndpoint
		}
		secureURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(baseEndpoint, "/"), s.bucketName, fileName)
	}

	return &UploadResult{
		PublicID:     fileName,
		SecureURL:    secureURL,
		Format:       strings.TrimPrefix(ext, "."),
		OriginalName: header.Filename,
		Bytes:        len(fileBytes),
	}, nil
}

func (s *r2StorageService) UploadData(data []byte, fileName string, contentType string, folder ImageFolder) (*UploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fullPath := fmt.Sprintf("%s/%s", folder, fileName)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fullPath),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("R2 upload failed: %w", err)
	}

	var secureURL string
	if s.publicURL != "" {
		secureURL = fmt.Sprintf("%s/%s", s.publicURL, fullPath)
	} else {
		baseEndpoint := ""
		if s.client.Options().BaseEndpoint != nil {
			baseEndpoint = *s.client.Options().BaseEndpoint
		}
		secureURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(baseEndpoint, "/"), s.bucketName, fullPath)
	}

	return &UploadResult{
		PublicID:     fullPath,
		SecureURL:    secureURL,
		Format:       strings.TrimPrefix(filepath.Ext(fileName), "."),
		OriginalName: fileName,
		Bytes:        len(data),
	}, nil
}

func (s *r2StorageService) DeleteImage(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(publicID),
	})

	if err != nil {
		return fmt.Errorf("failed to delete image from R2: %w", err)
	}
	return nil
}
