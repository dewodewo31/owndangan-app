package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name        string
	Size        int64
	ContentType string
	Extension   string
}

type UploadResult struct {
	URL      string
	Key      string
	Size     int64
	MIMEType string
}

type Storage interface {
	Upload(ctx context.Context, key string, data io.Reader, opts UploadOptions) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) string
	Exists(ctx context.Context, key string) (bool, error)
}

type UploadOptions struct {
	ContentType string
	Extension   string
	MaxSize     int64
	Metadata    map[string]string
}

type LocalStorage struct {
	basePath  string
	publicURL string
}

func NewLocalStorage(basePath, publicURL string) *LocalStorage {
	return &LocalStorage{
		basePath:  basePath,
		publicURL: publicURL,
	}
}

func (s *LocalStorage) Upload(ctx context.Context, key string, data io.Reader, opts UploadOptions) (*UploadResult, error) {
	if s.basePath != "" {
		fullPath := filepath.Join(s.basePath, filepath.FromSlash(key))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return nil, err
		}
		dst, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, data); err != nil {
			return nil, err
		}
	}
	result := &UploadResult{
		URL:      s.url(key),
		Key:      key,
		MIMEType: opts.ContentType,
	}
	return result, nil
}

func (s *LocalStorage) url(key string) string {
	if s.publicURL == "" {
		return "/" + key
	}
	return s.publicURL + "/" + key
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if s.basePath == "" {
		return nil
	}
	return os.Remove(filepath.Join(s.basePath, filepath.FromSlash(key)))
}

func (s *LocalStorage) GetURL(ctx context.Context, key string) string {
	return s.url(key)
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	if s.basePath == "" {
		return false, nil
	}
	_, err := os.Stat(filepath.Join(s.basePath, filepath.FromSlash(key)))
	if err == nil {
		return true, nil
	}
	return false, nil
}

type S3Storage struct {
	bucket    string
	region    string
	accessKey string
	secretKey string
	endpoint  string
	publicURL string
}

func NewS3Storage(bucket, region, accessKey, secretKey, endpoint, publicURL string) *S3Storage {
	return &S3Storage{
		bucket:    bucket,
		region:    region,
		accessKey: accessKey,
		secretKey: secretKey,
		endpoint:  endpoint,
		publicURL: publicURL,
	}
}

func (s *S3Storage) Upload(ctx context.Context, key string, data io.Reader, opts UploadOptions) (*UploadResult, error) {
	result := &UploadResult{
		URL:      s.publicURL + "/" + key,
		Key:      key,
		MIMEType: opts.ContentType,
	}
	return result, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *S3Storage) GetURL(ctx context.Context, key string) string {
	return s.publicURL + "/" + key
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func GenerateKey(prefix, ext string) string {
	return prefix + "/" + generateUUID() + ext
}

func generateUUID() string {
	return time.Now().Format("20060102") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
