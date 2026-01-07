package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// UploadStats represent the metadata of uploading file
type UploadStats struct {
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// In-memory store for upload stats
var uploads []UploadStats

//go:embed upload.html
var htmlUploadForm string

// Maximum file size (10 MB)
const maxFileSize = 10 * 1024 * 1024

func main() {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal(err)
	}

	// Create Echo instance
	e := echo.New()

	// Serve static files from the uploads directory
	e.Static("/files", "./uploads")

	// Serve the HTML upload form
	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return c.String(http.StatusOK, htmlUploadForm)
	})

	// Handle file upload
	e.POST("/upload", func(c echo.Context) error {
		// Get the uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Check file size
		if file.Size > maxFileSize {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf("File too large (max %d MB)", maxFileSize/(1024*1024)),
			)
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		// Create a safe filename
		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		basename := strings.TrimSuffix(filename, ext)
		filename = fmt.Sprintf("%s_%d%s", basename, time.Now().Unix(), ext)

		// Create destination file
		dst, err := os.Create(filepath.Join("uploads", filename))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer dst.Close()

		// Copy file contents
		if _, err := io.Copy(dst, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Detect MIME type
		mimeType := file.Header.Get("Content-Type")

		// Store upload stats
		stats := UploadStats{
			Filename:   filename,
			Size:       file.Size,
			MimeType:   mimeType,
			UploadedAt: time.Now(),
		}
		uploads = append(uploads, stats)

		return c.JSON(http.StatusOK, stats)
	})

	// Get list of uploaded files
	e.GET("/files-list", func(c echo.Context) error {
		return c.JSON(http.StatusOK, uploads)
	})

	// Start server
	log.Println("Starting file upload server on :8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
