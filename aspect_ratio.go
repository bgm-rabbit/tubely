package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type ffprobeStream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ffprobeOutput struct {
	Streams []ffprobeStream `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	// Run ffprobe command
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	// Capture stdout to a buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffprobe command failed: %w", err)
	}

	// Unmarshal the JSON output
	var output ffprobeOutput
	err = json.Unmarshal(stdout.Bytes(), &output)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal ffprobe output: %w", err)
	}

	if len(output.Streams) == 0 {
		return "", fmt.Errorf("no streams found in video")
	}

	// Get width and height from the first stream
	width := output.Streams[0].Width
	height := output.Streams[0].Height

	if width == 0 || height == 0 {
		return "other", nil
	}

	// Calculate aspect ratio using integer division for tolerance
	// 16:9 ratio has width/height ≈ 1.78
	// 9:16 ratio has width/height ≈ 0.56
	ratio := float64(width) / float64(height)

	log.Printf("Video dimensions: %dx%d, ratio: %.2f", width, height, ratio)

	// Use tolerance ranges
	if ratio > 1.4 { // Closer to 16:9 (1.78)
		return "16:9", nil
	} else if ratio < 0.8 { // Closer to 9:16 (0.56)
		return "9:16", nil
	}

	return "other", nil
}
