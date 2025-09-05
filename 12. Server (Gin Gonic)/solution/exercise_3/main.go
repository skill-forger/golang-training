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

	"github.com/gin-gonic/gin"
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
	err := os.MkdirAll("./uploads", 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Gin router with default middleware
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Serve static files from the uploads directory
	r.Static("/files", "./uploads")

	// Serve the HTML upload form
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlUploadForm)
	})

	// Handle file upload
	r.POST("/upload", func(c *gin.Context) {
		// Get the file from the request
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Check file size
		if header.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("File too large (max %d MB)", maxFileSize/(1024*1024)),
			})
			return
		}

		// Create a safe filename
		filename := header.Filename
		// Remove any path from the filename
		filename = filepath.Base(filename)
		// Ensure filename is unique
		ext := filepath.Ext(filename)
		basename := strings.TrimSuffix(filename, ext)
		filename = fmt.Sprintf("%s_%d%s", basename, time.Now().Unix(), ext)

		// Create destination file
		dst, err := os.Create(filepath.Join("uploads", filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer dst.Close()

		// Copy the file
		_, err = io.Copy(dst, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Detect mime type
		mimeType := header.Header.Get("Content-Type")

		// Store upload stats
		stats := UploadStats{
			Filename:   filename,
			Size:       header.Size,
			MimeType:   mimeType,
			UploadedAt: time.Now(),
		}
		uploads = append(uploads, stats)

		c.JSON(http.StatusOK, stats)
	})

	// Get list of uploaded files
	r.GET("/files-list", func(c *gin.Context) {
		c.JSON(http.StatusOK, uploads)
	})

	// Start the server
	log.Println("Starting file upload server on :8080...")
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
