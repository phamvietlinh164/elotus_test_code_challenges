package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const MaxImageSize = 8 << 20 // 8 MB

type Service interface {
	UploadImage(userID uint, file multipart.File, header *multipart.FileHeader, contentType, ip, userAgent string) (*Image, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) UploadImage(userID uint, file multipart.File, header *multipart.FileHeader, contentType, ip, userAgent string) (*Image, error) {
	// Size limit
	if header.Size > MaxImageSize {
		return nil, errors.New("file too large (max 8MB)")
	}

	// Check content type
	if contentType == "" || contentType[:5] != "image" {
		return nil, errors.New("invalid content type, must be image")
	}

	// Save to /tmp
	tmpFileName := fmt.Sprintf("/tmp/%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	out, err := os.Create(tmpFileName)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, err
	}

	meta := &Image{
		UserID:       userID,
		FileName:     tmpFileName,
		ContentType:  contentType,
		Size:         header.Size,
		OriginalName: header.Filename,
		UploadIP:     ip,
		UserAgent:    userAgent,
		UploadedAt:   time.Now().UTC(),
	}

	if err := s.repo.Save(meta); err != nil {
		return nil, err
	}

	return meta, nil
}
