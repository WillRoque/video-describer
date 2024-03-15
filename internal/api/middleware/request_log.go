package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/WillRoque/ai-video-descriptor/internal/api/context"
	"github.com/rs/zerolog/log"
)

type logEntry struct {
	RequestID     string
	ReceivedTime  time.Time
	RequestMethod string
	RequestURL    string
	UserAgent     string
	Referer       string
	Proto         string

	RemoteIP string
	ServerIP string

	Latency time.Duration
}

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		le := &logEntry{
			RequestID:     context.RequestID(r.Context()),
			ReceivedTime:  start,
			RequestMethod: r.Method,
			RequestURL:    r.URL.String(),
			UserAgent:     r.UserAgent(),
			Referer:       r.Referer(),
			Proto:         r.Proto,
			RemoteIP:      ipFromHostPort(r.RemoteAddr),
		}

		if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
			le.ServerIP = ipFromHostPort(addr.String())
		}

		next.ServeHTTP(w, r)

		le.Latency = time.Since(start)
		log.Info().
			Str("request_id", le.RequestID).
			Time("received_time", le.ReceivedTime).
			Str("method", le.RequestMethod).
			Str("url", le.RequestURL).
			Str("agent", le.UserAgent).
			Str("referer", le.Referer).
			Str("proto", le.Proto).
			Str("remote_ip", le.RemoteIP).
			Str("server_ip", le.ServerIP).
			Dur("latency", le.Latency).
			Msg("")
	})
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
