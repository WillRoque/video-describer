package youtube

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the YoutubeClient interface.
type mockYoutubeClient struct {
	mock.Mock
}

func (m *mockYoutubeClient) GetVideo(id string) (*youtube.Video, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*youtube.Video), nil
	}
	return nil, args.Error(1)
}

func (m *mockYoutubeClient) GetStream(video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error) {
	args := m.Called(video, format)
	return args.Get(0).(io.ReadCloser), args.Get(1).(int64), args.Error(2)
}

// Mock io.ReadCloser.
type mockStream struct{}

func (s mockStream) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (s mockStream) Close() error {
	return nil
}

func TestDownloadVideo(t *testing.T) {
	videoID := "dQw4w9WgXcQ"
	mockClient := &mockYoutubeClient{}
	expectedVideo := &youtube.Video{
		ID:      videoID,
		Formats: youtube.FormatList{youtube.Format{Quality: "360p"}},
	}

	mockClient.
		On("GetVideo", videoID).
		Return(expectedVideo, nil)
	mockClient.
		On("GetStream", expectedVideo, mock.AnythingOfType("*youtube.Format")).
		Return(mockStream{}, int64(1024), nil)

	ytDownloader := NewDownloader(mockClient)
	videoPath, err := ytDownloader.DownloadVideo(videoID)
	assert.NoError(t, err)
	assert.FileExists(t, videoPath)

	cleanup(t, videoPath)
}

func cleanup(t *testing.T, videoPath string) {
	videoDir := filepath.Dir(videoPath)
	err := os.RemoveAll(videoDir)
	assert.NoError(t, err)
}
