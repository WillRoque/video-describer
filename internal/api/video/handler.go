package video

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	c "github.com/WillRoque/ai-video-descriptor/internal/api/context"
	e "github.com/WillRoque/ai-video-descriptor/internal/api/error"
	"github.com/WillRoque/ai-video-descriptor/internal/api/logging"
	"github.com/WillRoque/ai-video-descriptor/internal/video"
	"github.com/rs/zerolog/log"
)

type VisionAI interface {
	DescribeVideo(ctx context.Context, framesPath []string, prompt string) (string, error)
}

type YouTubeDownloader interface {
	DownloadVideo(videoID string) (string, error)
}

type API struct {
	visionAI          VisionAI
	youtubeDownloader YouTubeDownloader
}

func New(visionAI VisionAI, youtubeDownloader YouTubeDownloader) *API {
	return &API{visionAI, youtubeDownloader}
}

func (a *API) Describe(w http.ResponseWriter, r *http.Request) {
	reqID := c.RequestID(r.Context())

	logging.LogHandlerInfo(r.Context(), "Decoding JSON payload")
	var dr DescribeRequest
	err := json.NewDecoder(r.Body).Decode(&dr)
	if err != nil {
		e.BadRequest(w, r, "json decode failure", err)
		return
	}

	logging.LogHandlerInfo(r.Context(), "Validating request")
	valErrs := dr.Validate()
	if valErrs != "" {
		e.ValidationError(w, r, valErrs)
		return
	}

	logging.LogHandlerInfo(r.Context(), "Downloading video")
	videoPath, err := a.youtubeDownloader.DownloadVideo(dr.VideoID)
	if err != nil {
		e.ServerError(w, r, err)
		return
	}
	defer cleanupVideo(reqID, videoPath)

	logging.LogHandlerInfo(r.Context(), "Extracting video frames")
	framesPaths, err := video.ExtractFrames(videoPath, dr.Start.Duration, dr.End.Duration)
	if err != nil {
		e.ServerError(w, r, err)
		return
	}

	logging.LogHandlerInfo(r.Context(), "Describing video")
	description, err := a.visionAI.DescribeVideo(r.Context(), framesPaths, dr.Prompt)
	if err != nil {
		e.ServerError(w, r, err)
		return
	}

	w.Write([]byte(description))
}

func cleanupVideo(reqID, videoPath string) {
	videoDir := filepath.Dir(videoPath)
	err := os.RemoveAll(videoDir)
	if err != nil {
		log.Error().Str("request_id", reqID).Msgf("failed to cleanup directory %s", videoDir)
	}
}
