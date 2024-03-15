package error

import (
	"net/http"

	c "github.com/WillRoque/ai-video-descriptor/internal/api/context"
	"github.com/rs/zerolog/log"
)

func BadRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	reqID := c.RequestID(r.Context())
	log.Error().Str(string(c.KeyRequestID), reqID).Err(err).Msg(msg)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func ValidationError(w http.ResponseWriter, r *http.Request, msg string) {
	reqID := c.RequestID(r.Context())
	log.Error().Str(string(c.KeyRequestID), reqID).Msg(msg)
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(msg))
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	reqID := c.RequestID(r.Context())
	log.Error().Str(string(c.KeyRequestID), reqID).Err(err).Msg("")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}
