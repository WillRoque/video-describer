package context

import "context"

const KeyRequestID key = "requestID"

type key string

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(KeyRequestID).(string)
	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID)
}
