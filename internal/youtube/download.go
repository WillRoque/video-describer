package youtube

import (
	"fmt"
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
)

type YoutubeClient interface {
	GetVideo(videoID string) (*youtube.Video, error)
	GetStream(video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error)
}

type Downloader struct {
	client YoutubeClient
}

func NewDownloader(client YoutubeClient) Downloader {
	return Downloader{client}
}

func (d *Downloader) DownloadVideo(videoID string) (string, error) {
	outDir, err := os.MkdirTemp("", videoID)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	videoPath := fmt.Sprintf("%s/video.mp4", outDir)

	video, err := d.client.GetVideo(videoID)
	if err != nil {
		return "", fmt.Errorf("failed to get video with ID %s: %w", videoID, err)
	}

	formats := video.Formats.Quality("360p")
	stream, _, err := d.client.GetStream(video, &formats[0])
	if err != nil {
		return "", fmt.Errorf("failed to get stream for video with ID %s: %w", videoID, err)
	}
	defer stream.Close()

	file, err := os.Create(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file at path %s: %w", videoPath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to copy stream to file: %w", err)
	}

	return videoPath, nil
}
