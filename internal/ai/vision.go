package ai

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

const defaultPrompt = "These are frames of a video. Describe what happens in the video in one sentence."

type Vision struct {
	client *openai.Client
}

func NewVision(client *openai.Client) Vision {
	return Vision{client}
}

func (v *Vision) DescribeVideo(ctx context.Context, framesPath []string, prompt string) (string, error) {
	if prompt == "" {
		prompt = defaultPrompt
	}

	content := []openai.ChatMessagePart{{
		Type: openai.ChatMessagePartTypeText,
		Text: prompt,
	}}

	framesMsgs, err := appendVideoFrames(framesPath)
	if err != nil {
		return "", fmt.Errorf("failed to append video frames: %w", err)
	}
	content = append(content, framesMsgs...)

	messages := []openai.ChatCompletionMessage{{
		Role:         openai.ChatMessageRoleUser,
		MultiContent: content,
	}}

	request := openai.ChatCompletionRequest{
		Model:    "gpt-4-1106-vision-preview",
		Messages: messages,
		Stream:   true,
	}

	return v.handleStreamRequest(ctx, request)
}

func (v *Vision) handleStreamRequest(ctx context.Context, request openai.ChatCompletionRequest) (string, error) {
	var out string
	stream, err := v.client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion stream: %w", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error receiving stream response: %w", err)
		}

		out += response.Choices[0].Delta.Content
	}

	return out, nil
}

func appendVideoFrames(framesPath []string) ([]openai.ChatMessagePart, error) {
	var content []openai.ChatMessagePart

	// Limit to 10 images because of a limitation with the API from OpenAI.
	maxEl := 10
	gap := max(len(framesPath)/maxEl, 1)

	for i := 0; i < maxEl && i < len(framesPath); i++ {
		frame := i * gap
		file, err := os.Open(framesPath[frame])
		if err != nil {
			return nil, fmt.Errorf("failed to open frame file %s: %w", framesPath[frame], err)
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read frame file %s: %w", framesPath[frame], err)
		}

		b64 := base64.StdEncoding.EncodeToString(fileContent)

		frameReq := openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{
				URL:    fmt.Sprintf("data:image/jpeg;base64,%s", b64),
				Detail: openai.ImageURLDetailLow,
			},
		}

		content = append(content, frameReq)
	}

	return content, nil
}
