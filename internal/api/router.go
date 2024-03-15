package api

import (
	"net/http"

	"github.com/WillRoque/ai-video-descriptor/internal/api/middleware"
	"github.com/WillRoque/ai-video-descriptor/internal/api/video"
	"github.com/go-chi/chi"
)

func NewRouter(visionAI video.VisionAI, ytDownloader video.YouTubeDownloader) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)
		r.Use(middleware.RequestLog)

		videoAPI := video.New(visionAI, ytDownloader)
		r.MethodFunc(http.MethodPost, "/video/describe", videoAPI.Describe)
	})

	return r
}
