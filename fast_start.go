package main

import (
	"fmt"
	"log"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	// Create output file path by appending .processing to the input file
	outputPath := filePath + ".processing"

	// Create ffmpeg command with fast start encoding
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg processing failed: %w", err)
	}

	log.Printf("Video processed for fast start: %s -> %s", filePath, outputPath)

	// Return the output file path
	return outputPath, nil
}
