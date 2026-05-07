package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	// Return video as-is if no VideoURL (e.g., before upload)
	if video.VideoURL == nil || *video.VideoURL == "" {
		return video, nil
	}

	// Split the video URL on comma to get bucket and key
	parts := strings.Split(*video.VideoURL, ",")
	if len(parts) != 2 {
		return video, fmt.Errorf("invalid video URL format: expected 'bucket,key'")
	}

	bucket := parts[0]
	key := parts[1]

	// Generate presigned URL with 24-hour expiration
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, 24*time.Hour)
	if err != nil {
		return video, fmt.Errorf("couldn't generate presigned URL: %w", err)
	}

	// Update the VideoURL field with the presigned URL
	video.VideoURL = &presignedURL

	return video, nil
}
