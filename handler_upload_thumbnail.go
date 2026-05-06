package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// Parse the form data
	const maxMemory = 10 << 20
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse form", err)
		return
	}

	// Get the image data from the form
	file, fileHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get form file", err)
		return
	}
	defer file.Close()

	// Get the media type from the Content-Type header
	contentType := fileHeader.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse media type", err)
		return
	}

	// Validate that the media type is either image/jpeg or image/png
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Only JPEG and PNG images are allowed", nil)
		return
	}

	// Get the video's metadata from the database
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get video", err)
		return
	}

	// Check if the authenticated user is the video owner
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You are not the owner of this video", nil)
		return
	}

	// Determine file extension from Content-Type
	ext := getFileExtension(mediaType)

	// Generate a random filename
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate random filename", err)
		return
	}
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)
	fileName := randomString + ext
	filePath := filepath.Join(cfg.assetsRoot, fileName)

	// Create the file on disk
	newFile, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create file", err)
		return
	}
	defer newFile.Close()

	// Copy the file contents
	_, err = io.Copy(newFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file", err)
		return
	}

	// Update the thumbnail URL
	thumbnailURL := fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, fileName)
	video.ThumbnailURL = &thumbnailURL

	// Update the video in the database
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	// Respond with the updated video JSON
	respondWithJSON(w, http.StatusOK, video)
}

func getFileExtension(mediaType string) string {
	switch mediaType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		// Try to extract from the media type if it has a subtype
		parts := strings.Split(mediaType, "/")
		if len(parts) > 1 {
			return "." + parts[1]
		}
		return ".bin"
	}
}
