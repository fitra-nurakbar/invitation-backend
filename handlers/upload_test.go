package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupUploadRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/upload", UploadFile)
	return r
}

func createMultipartRequest(t *testing.T, fieldName, filename string, content []byte) (*http.Request, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to copy content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestUploadFileValidation(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	router := setupUploadRouter()

	t.Run("Reject unsupported file extension", func(t *testing.T) {
		req, _ := createMultipartRequest(t, "file", "malicious.php", []byte("<?php echo 'hack'; ?>"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for .php file, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Reject SVG upload due to stored XSS risk", func(t *testing.T) {
		req, _ := createMultipartRequest(t, "file", "test.svg", []byte("<svg><script>alert(1)</script></svg>"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for .svg file, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Reject spoofed JPEG with plain text payload", func(t *testing.T) {
		req, _ := createMultipartRequest(t, "file", "fake.jpg", []byte("plain text disguised as jpg"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for spoofed jpg, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Reject file exceeding size limit", func(t *testing.T) {
		// Image limit is 10MB; create 11MB spoofed payload
		oversized := make([]byte, 11<<20)
		oversized[0] = 0xff
		oversized[1] = 0xd8
		oversized[2] = 0xff
		req, _ := createMultipartRequest(t, "file", "large.jpg", oversized)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Fatalf("expected 413 for oversized file, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Accept valid PNG file", func(t *testing.T) {
		// Valid minimal 1x1 PNG bytes
		validPNG := []byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, // PNG header
			0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
			0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
			0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
			0x42, 0x60, 0x82,
		}
		req, _ := createMultipartRequest(t, "file", "test.png", validPNG)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 for valid PNG, got %d: %s", w.Code, w.Body.String())
		}
	})
}
