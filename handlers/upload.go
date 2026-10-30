package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Batas ukuran maksimal file: 15 MB.
const MaxUploadSize = 15 << 20
const maxImageSize = 10 << 20
const maxAudioSize = 15 << 20

var allowedUploadTypes = map[string]struct {
	mimeType string
	maxSize  int64
}{
	".jpg":  {"image/jpeg", maxImageSize},
	".jpeg": {"image/jpeg", maxImageSize},
	".png":  {"image/png", maxImageSize},
	".webp": {"image/webp", maxImageSize},
	".gif":  {"image/gif", maxImageSize},
	".mp3":  {"audio/mpeg", maxAudioSize},
	".wav":  {"audio/wav", maxAudioSize},
	".m4a":  {"audio/mp4", maxAudioSize},
	".aac":  {"audio/aac", maxAudioSize},
	".ogg":  {"audio/ogg", maxAudioSize},
}

// POST /api/v1/upload
func UploadFile(c *gin.Context) {
	// Limit total multipart body; the multipart overhead is allowed by this request cap.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize+(1<<20))

	fileHeader, err := c.FormFile("file")
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Ukuran request melebihi batas upload"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan atau request multipart tidak valid"})
		return
	}

	ext := strings.ToLower(filepath.Ext(filepath.Base(fileHeader.Filename)))
	fileType, isAllowed := allowedUploadTypes[ext]
	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung"})
		return
	}
	if fileHeader.Size <= 0 || fileHeader.Size > fileType.maxSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Ukuran file melebihi batas untuk format ini"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak dapat dibaca"})
		return
	}
	defer file.Close()

	// Inspect the actual bytes instead of trusting the filename or client MIME header.
	sniff := make([]byte, 512)
	n, readErr := io.ReadFull(file, sniff)
	if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca isi file"})
		return
	}
	sniff = sniff[:n]
	if !matchesUploadSignature(ext, sniff) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Isi file tidak sesuai dengan format yang diizinkan"})
		return
	}
	if contentType := http.DetectContentType(sniff); ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		mediaType, _, parseErr := mime.ParseMediaType(contentType)
		if parseErr != nil || mediaType != fileType.mimeType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe konten file tidak valid"})
			return
		}
	}

	// Create files with restrictive permissions and a server-generated, unguessable name.
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan direktori penyimpanan"})
		return
	}

	uniqueFilename := uuid.NewString() + ext
	savePath := filepath.Join(uploadDir, uniqueFilename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file di server"})
		return
	}
	if err := os.Chmod(savePath, 0640); err != nil {
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengatur izin file"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
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
			"content_type": fileType.mimeType,
		},
	})
}

func matchesUploadSignature(ext string, data []byte) bool {
	switch ext {
	case ".jpg", ".jpeg":
		return len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff
	case ".png":
		return len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n"
	case ".gif":
		return len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a")
	case ".webp":
		return len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP"
	case ".mp3":
		return len(data) >= 3 && string(data[:3]) == "ID3" || len(data) >= 2 && data[0] == 0xff && data[1]&0xe0 == 0xe0
	case ".wav":
		return len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WAVE"
	case ".m4a":
		return len(data) >= 12 && string(data[4:8]) == "ftyp"
	case ".aac":
		return len(data) >= 2 && data[0] == 0xff && data[1]&0xf6 == 0xf0
	case ".ogg":
		return len(data) >= 4 && string(data[:4]) == "OggS"
	default:
		return false
	}
}
