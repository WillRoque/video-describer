package video

import (
	"fmt"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// ExtractFrames extracts one frame per second of the entire video,
// or from a specific start and end.
func ExtractFrames(videoPath string, start, end time.Duration) ([]string, error) {
	videoDir := filepath.Dir(videoPath)
	outFiles := fmt.Sprintf("%s/frame%%03d.jpg", videoDir)

	inArgs := ffmpeg.KwArgs{}
	if start != 0 {
		inArgs["ss"] = start.Seconds()
	}
	if end != 0 {
		inArgs["to"] = end.Seconds()
	}

	err := ffmpeg.Input(videoPath, inArgs).
		Output(outFiles, ffmpeg.KwArgs{"r": "1"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return []string{}, fmt.Errorf("failed to run ffmpeg command: %w", err)
	}

	var framesPaths []string
	framesPaths, err = filepath.Glob(filepath.Join(videoDir, "frame*"))
	if err != nil {
		return []string{}, fmt.Errorf("failed to glob frames: %w", err)
	}

	if len(framesPaths) == 0 {
		return []string{}, fmt.Errorf("could not extract frames from the video, check start and end parameters")
	}

	return framesPaths, nil
}
