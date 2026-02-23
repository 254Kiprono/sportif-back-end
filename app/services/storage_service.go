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
	"sync"
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

// b2StorageService uses the native Backblaze B2 HTTP API directly.
// Authorization is lazy — it happens on first upload, not at startup.
type b2StorageService struct {
	mu             sync.Mutex
	keyID          string
	applicationKey string
	bucketName     string
	// Populated after authorization
	bucketID    string
	downloadURL string
	apiURL      string
	authToken   string
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	if cfg.B2KeyID == "" || cfg.B2ApplicationKey == "" || cfg.B2BucketName == "" {
		return nil, fmt.Errorf("B2_KEY_ID, B2_APPLICATION_KEY and B2_BUCKET_NAME must be set")
	}

	return &b2StorageService{
		keyID:          strings.TrimSpace(cfg.B2KeyID),
		applicationKey: strings.TrimSpace(cfg.B2ApplicationKey),
		bucketName:     strings.TrimSpace(cfg.B2BucketName),
	}, nil
}

// ensureAuthorized calls b2_authorize_account if not yet authorized (lazy init).
func (s *b2StorageService) ensureAuthorized() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.authToken != "" && s.apiURL != "" {
		return nil // already authorized
	}

	creds := base64.StdEncoding.EncodeToString([]byte(s.keyID + ":" + s.applicationKey))
	req, err := http.NewRequest("GET", "https://api.backblazeb2.com/b2api/v2/b2_authorize_account", nil)
	if err != nil {
		return fmt.Errorf("could not build auth request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+creds)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("B2 auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("B2 auth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp struct {
		AuthorizationToken string `json:"authorizationToken"`
		APIURL             string `json:"apiUrl"`
		DownloadURL        string `json:"downloadUrl"`
	}
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("could not decode B2 auth response: %w", err)
	}

	s.authToken = authResp.AuthorizationToken
	s.apiURL = authResp.APIURL
	s.downloadURL = authResp.DownloadURL
	return nil
}

// reauthorize forces a fresh authorization (call when token expires)
func (s *b2StorageService) reauthorize() error {
	s.mu.Lock()
	s.authToken = ""
	s.bucketID = ""
	s.mu.Unlock()
	return s.ensureAuthorized()
}

// getBucketID resolves the bucket name to its B2 bucket ID.
func (s *b2StorageService) getBucketID(ctx context.Context) (string, error) {
	s.mu.Lock()
	if s.bucketID != "" {
		id := s.bucketID
		s.mu.Unlock()
		return id, nil
	}
	s.mu.Unlock()

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
		return "", fmt.Errorf("b2_list_buckets request failed: %w", err)
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
			s.mu.Lock()
			s.bucketID = b.BucketID
			s.mu.Unlock()
			return s.bucketID, nil
		}
	}
	return "", fmt.Errorf("bucket '%s' not found in B2 account", s.bucketName)
}

// getUploadURL calls b2_get_upload_url to get a one-time upload URL.
func (s *b2StorageService) getUploadURL(ctx context.Context) (string, string, error) {
	bucketID, err := s.getBucketID(ctx)
	if err != nil {
		return "", "", err
	}

	body, _ := json.Marshal(map[string]string{"bucketId": bucketID})
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_get_upload_url", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", s.authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("b2_get_upload_url failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("b2_get_upload_url failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		UploadURL          string `json:"uploadUrl"`
		AuthorizationToken string `json:"authorizationToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	return result.UploadURL, result.AuthorizationToken, nil
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
		return nil, fmt.Errorf("file size exceeds the 10MB limit")
	}

	// Lazy authorization
	if err := s.ensureAuthorized(); err != nil {
		return nil, err
	}

	// Read entire file into memory (required for SHA1 and exact Content-Length)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read uploaded file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(fileBytes)

	// Compute SHA1 (required by B2)
	h := sha1.New()
	h.Write(fileBytes)
	sha1Hex := fmt.Sprintf("%x", h.Sum(nil))

	// Build unique file path
	timestamp := time.Now().UnixMilli()
	cleanName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	fileName := fmt.Sprintf("%s/%s_%d%s", folder, cleanName, timestamp, ext)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get a one-time upload URL
	uploadURL, uploadToken, err := s.getUploadURL(ctx)
	if err != nil {
		// Token may have expired — re-authorize and retry once
		if err2 := s.reauthorize(); err2 == nil {
			uploadURL, uploadToken, err = s.getUploadURL(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("could not get B2 upload URL: %w", err)
		}
	}

	// Upload the file
	uploadReq, err := http.NewRequestWithContext(ctx, "POST", uploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("could not build upload request: %w", err)
	}
	uploadReq.ContentLength = int64(len(fileBytes))
	uploadReq.Header.Set("Authorization", uploadToken)
	uploadReq.Header.Set("X-Bz-File-Name", fileName)
	uploadReq.Header.Set("Content-Type", contentType)
	uploadReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	uploadReq.Header.Set("X-Bz-Content-Sha1", sha1Hex)

	httpClient := &http.Client{Timeout: 60 * time.Second}
	uploadResp, err := httpClient.Do(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("B2 upload request failed: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("B2 upload failed (status %d): %s", uploadResp.StatusCode, string(respBody))
	}

	// Construct public download URL
	secureURL := fmt.Sprintf("%s/file/%s/%s",
		strings.TrimSuffix(s.downloadURL, "/"),
		s.bucketName,
		fileName,
	)

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
	if err := s.ensureAuthorized(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketID, err := s.getBucketID(ctx)
	if err != nil {
		return err
	}

	// List to get the fileId
	listBody, _ := json.Marshal(map[string]interface{}{
		"bucketId":     bucketID,
		"prefix":       publicID,
		"maxFileCount": 1,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL+"/b2api/v2/b2_list_file_names", bytes.NewReader(listBody))
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
		return nil // already deleted or doesn't exist
	}

	deleteBody, _ := json.Marshal(map[string]string{
		"fileId":   listResult.Files[0].FileID,
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
	delResp.Body.Close()
	return nil
}
