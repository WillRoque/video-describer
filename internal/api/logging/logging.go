package logging

import (
	"context"

	c "github.com/WillRoque/ai-video-descriptor/internal/api/context"
	"github.com/rs/zerolog/log"
)

func LogHandlerInfo(ctx context.Context, msg string) {
	reqID := c.RequestID(ctx)
	log.Info().Str(string(c.KeyRequestID), reqID).Msg(msg)
}
