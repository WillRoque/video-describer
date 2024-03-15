package ai

import (
	"context"
	"net/http"
	"runtime"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

var framePath string

func init() {
	// Get the path to this test file as if it were the path to a video frame.
	// For the purpose of the test, any file works.
	_, framePath, _, _ = runtime.Caller(0)
}

func setupOpenAITestServer() (client *openai.Client, server *ServerTest, teardown func()) {
	server = NewTestServer()
	ts := server.OpenAITestServer()
	ts.Start()
	teardown = ts.Close
	config := openai.DefaultConfig(GetTestToken())
	config.BaseURL = ts.URL + "/v1"
	client = openai.NewClientWithConfig(config)
	return
}

func TestDescribeVideo(t *testing.T) {
	client, server, teardown := setupOpenAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		// Test responses
		dataBytes := []byte{}
		dataBytes = append(dataBytes, []byte("event: message\n")...)

		data := `{"id":"1","object":"completion","created":1598069254,"model":"gpt-4-1106-vision-preview","choices":[{"index":0,"delta":{"content":"response1"},"finish_reason":"max_tokens"}]}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: message\n")...)

		data = `{"id":"2","object":"completion","created":1598069255,"model":"gpt-4-1106-vision-preview","choices":[{"index":0,"delta":{"content":"response2"},"finish_reason":"max_tokens"}]}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: done\n")...)
		dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

		_, err := w.Write(dataBytes)
		assert.NoError(t, err)
	})

	vision := NewVision(client)
	framesPath := []string{framePath}

	description, err := vision.DescribeVideo(context.Background(), framesPath, "")
	assert.NoError(t, err)
	assert.Greater(t, len(description), 0)
}

func TestAppendVideoFrames(t *testing.T) {
	framesPath := []string{framePath, framePath}

	parts, err := appendVideoFrames(framesPath)
	assert.NoError(t, err)
	assert.Len(t, parts, 2)
}

func TestAppendVideoFramesInvalidPath(t *testing.T) {
	framesPath := []string{"invalid/path.jpg"}

	_, err := appendVideoFrames(framesPath)
	assert.Error(t, err)
}

func TestAppendVideoFramesMoreThan10(t *testing.T) {
	framesPath := make([]string, 100)
	for i := 0; i < 100; i++ {
		framesPath[i] = framePath
	}

	parts, err := appendVideoFrames(framesPath)
	assert.NoError(t, err)
	assert.Len(t, parts, 10)
}
