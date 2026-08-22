package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	jwtoken "github.com/riazahmedshah/vfs/internal/model/jwt"
)

func GenerateJWT(userID uuid.UUID, storageLimit int64, FileSizeLimit int64, isGuest bool, expiry int, JWTSecret string) (string, error) {
	claims := &jwtoken.CustomClaims{
		UserID:        userID,
		StorageLimit:  storageLimit,
		FileSizeLimit: FileSizeLimit,
		IsGuest:       isGuest,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

var allowedMimeToExt = map[string][]string{
	"image/jpeg":         {"jpg", "jpeg"},
	"image/png":          {"png"},
	"application/pdf":    {"pdf"},
	"text/plain":         {"txt"},
	"video/mp4":          {"mp4"},
	"application/zip":    {"zip"},
	"application/msword": {"doc"},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {"docx"},
	"application/vnd.ms-excel": {"xls"},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         {"xlsx"},
	"application/vnd.ms-powerpoint":                                             {"ppt"},
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": {"pptx"},
}

func DetectAndValidateExt(src multipart.File, claimedExt string) (string, error) {
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	// Read ke baad file pointer wapas start pe le jao, warna GCS upload me content miss ho jayega
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	contentType := http.DetectContentType(buf[:n])

	allowedExts, ok := allowedMimeToExt[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	if !slices.Contains(allowedExts, strings.ToLower(claimedExt)) {
		return "", fmt.Errorf("file extension %q does not match detected content type %s", claimedExt, contentType)
	}

	return contentType, nil
}
