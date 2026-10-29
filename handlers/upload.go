package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Batas ukuran maksimal file: 15 MB
const MaxUploadSize = 15 << 20

// POST /api/v1/upload
func UploadFile(c *gin.Context) {
	// Batasi ukuran request body
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File tidak ditemukan atau ukuran melebihi batas 15MB",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtensions := map[string]string{
		// Format gambar
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		// Format audio
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".m4a":  "audio/m4a",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
	}

	expectedMime, isAllowed := allowedExtensions[ext]
	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format file tidak didukung. Format yang diizinkan: JPG, PNG, WEBP, GIF, SVG, MP3, WAV, M4A, AAC, OGG",
		})
		return
	}

	// Buat folder ./uploads jika belum tersedia
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyiapkan direktori penyimpanan",
		})
		return
	}

	// Generate nama file unik aman dari path traversal
	randomID := uuid.New().String()[:8]
	uniqueFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomID, ext)
	savePath := filepath.Join(uploadDir, uniqueFilename)

	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan file di server",
		})
		return
	}

	// Tentukan skema URL
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	// Jika ada BASE_URL di env gunakan itu, jika tidak gunakan host request
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}
	baseURL = strings.TrimRight(baseURL, "/")

	fileURL := fmt.Sprintf("%s/uploads/%s", baseURL, uniqueFilename)

	c.JSON(http.StatusOK, gin.H{
		"message": "File berhasil diunggah",
		"data": gin.H{
			"url":          fileURL,
			"path":         fmt.Sprintf("/uploads/%s", uniqueFilename),
			"filename":     uniqueFilename,
			"size":         fileHeader.Size,
			"content_type": expectedMime,
		},
	})
}
