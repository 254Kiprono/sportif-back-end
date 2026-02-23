package services

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"webuye-sportif/app/config"
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

// b2StorageService uses the native Backblaze B2 HTTP API directly,
// bypassing the AWS S3 SDK to avoid signature compatibility issues.
type b2StorageService struct {
	keyID          string
	applicationKey string
	bucketName     string
	bucketID       string
	downloadURL    string
	apiURL         string
	authToken      string
}

// b2AuthResponse is the response from b2_authorize_account
type b2AuthResponse struct {
	AuthorizationToken string `json:"authorizationToken"`
	APIURL             string `json:"apiUrl"`
	DownloadURL        string `json:"downloadUrl"`
	AllowedBucketName  string `json:"allowed"`
}

// b2UploadURLResponse is the response from b2_get_upload_url
type b2UploadURLResponse struct {
	UploadURL          string `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

// b2UploadFileResponse is the response from the actual file upload
type b2UploadFileResponse struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	if cfg.B2KeyID == "" || cfg.B2ApplicationKey == "" || cfg.B2BucketName == "" {
		return nil, fmt.Errorf("Backblaze B2 credentials are not fully configured")
	}

	keyID := strings.TrimSpace(cfg.B2KeyID)
	appKey := strings.TrimSpace(cfg.B2ApplicationKey)
	bucketName := strings.TrimSpace(cfg.B2BucketName)

	svc := &b2StorageService{
		keyID:          keyID,
		applicationKey: appKey,
		bucketName:     bucketName,
	}

	// Authorize immediately to get API URL and auth token
	if err := svc.authorize(); err != nil {
		return nil, fmt.Errorf("B2 authorization failed: %w", err)
	}

	return svc, nil
}

// authorize calls b2_authorize_account to get the API URL and token
func (s *b2StorageService) authorize() error {
	creds := base64.StdEncoding.EncodeToString([]byte(s.keyID + ":" + s.applicationKey))

	req, err := http.NewRequest("GET", "https://api.backblazeb2.com/b2api/v2/b2_authorize_account", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+creds)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp struct {
		AuthorizationToken string `json:"authorizationToken"`
		APIURL             string `json:"apiUrl"`
		DownloadURL        string `json:"downloadUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("could not decode auth response: %w", err)
	}

	s.authToken = authResp.AuthorizationToken
	s.apiURL = authResp.APIURL
	s.downloadURL = authResp.DownloadURL
	return nil
}

// getUploadURL calls b2_get_upload_url to get a one-time upload URL for the bucket
func (s *b2StorageService) getUploadURL(ctx context.Context) (*b2UploadURLResponse, error) {
	// We need the bucket ID — look it up from the bucket name
	bucketID, err := s.getBucketID(ctx)
	if err != nil {
		return nil, err
	}

	body, _ := json.Marshal(map[string]string{"bucketId": bucketID})
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_get_upload_url", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get_upload_url request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get_upload_url failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result b2UploadURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// getBucketID resolves a bucket name to its ID
func (s *b2StorageService) getBucketID(ctx context.Context) (string, error) {
	if s.bucketID != "" {
		return s.bucketID, nil
	}

	body, _ := json.Marshal(map[string]string{"bucketName": s.bucketName})
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_list_buckets", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("list_buckets request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Buckets []struct {
			BucketID   string `json:"bucketId"`
			BucketName string `json:"bucketName"`
		} `json:"buckets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	for _, b := range result.Buckets {
		if b.BucketName == s.bucketName {
			s.bucketID = b.BucketID
			return s.bucketID, nil
		}
	}
	return "", fmt.Errorf("bucket '%s' not found", s.bucketName)
}

func (s *b2StorageService) UploadImage(
	file multipart.File,
	header *multipart.FileHeader,
	folder ImageFolder,
) (*UploadResult, error) {
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
		return nil, fmt.Errorf("file size exceeds 10MB limit")
	}

	// Read entire file into memory for reliable upload
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(fileBytes)

	// Compute SHA1 for B2 integrity verification
	h := sha1.New()
	h.Write(fileBytes)
	sha1Hex := fmt.Sprintf("%x", h.Sum(nil))

	// Build unique file name
	timestamp := time.Now().UnixMilli()
	cleanName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	// B2 native API uses forward slashes; URL-encode folder separators are fine
	fileName := fmt.Sprintf("%s/%s_%d%s", folder, cleanName, timestamp, ext)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get a one-time upload URL from B2
	uploadInfo, err := s.getUploadURL(ctx)
	if err != nil {
		// Token may have expired — re-authorize and retry once
		if err2 := s.authorize(); err2 == nil {
			uploadInfo, err = s.getUploadURL(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("could not get B2 upload URL: %w", err)
		}
	}

	// Upload the file to B2
	uploadReq, err := http.NewRequestWithContext(ctx, "POST", uploadInfo.UploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}
	uploadReq.Header.Set("Authorization", uploadInfo.AuthorizationToken)
	uploadReq.Header.Set("X-Bz-File-Name", fileName)
	uploadReq.Header.Set("Content-Type", contentType)
	uploadReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	uploadReq.Header.Set("X-Bz-Content-Sha1", sha1Hex)
	uploadReq.ContentLength = int64(len(fileBytes))

	httpClient := &http.Client{Timeout: 60 * time.Second}
	uploadResp, err := httpClient.Do(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("B2 upload request failed: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("B2 upload failed (status %d): %s", uploadResp.StatusCode, string(body))
	}

	// Construct public download URL
	// B2 public download URL: {downloadUrl}/file/{bucketName}/{fileName}
	secureURL := fmt.Sprintf("%s/file/%s/%s", strings.TrimSuffix(s.downloadURL, "/"), s.bucketName, fileName)

	return &UploadResult{
		PublicID:     fileName,
		SecureURL:    secureURL,
		Format:       strings.TrimPrefix(ext, "."),
		OriginalName: header.Filename,
		Bytes:        len(fileBytes),
	}, nil
}

// DeleteImage removes an image from B2 using the native API.
func (s *b2StorageService) DeleteImage(publicID string) error {
	// To delete, we need fileId first — list file versions and grab it
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketID, err := s.getBucketID(ctx)
	if err != nil {
		return err
	}

	// List file versions to get the fileId
	body, _ := json.Marshal(map[string]interface{}{
		"bucketId":     bucketID,
		"prefix":       publicID,
		"maxFileCount": 1,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_list_file_names", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var listResult struct {
		Files []struct {
			FileID   string `json:"fileId"`
			FileName string `json:"fileName"`
		} `json:"files"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResult); err != nil {
		return err
	}
	if len(listResult.Files) == 0 {
		return nil // File doesn't exist, nothing to delete
	}

	fileID := listResult.Files[0].FileID
	deleteBody, _ := json.Marshal(map[string]string{
		"fileId":   fileID,
		"fileName": publicID,
	})
	delReq, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_delete_file_version", bytes.NewReader(deleteBody))
	if err != nil {
		return err
	}
	delReq.Header.Set("Authorization", s.authToken)
	delReq.Header.Set("Content-Type", "application/json")

	delResp, err := client.Do(delReq)
	if err != nil {
		return err
	}
	defer delResp.Body.Close()
	return nil
}
