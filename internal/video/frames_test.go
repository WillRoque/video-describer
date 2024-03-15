package video

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testVideoPath is the path to the test video file located in the testdata directory.
var testVideoPath string

func init() {
	// This determines the path to the test video file at runtime.
	_, testFile, _, _ := runtime.Caller(0)
	testdataDir := filepath.Join(path.Dir(testFile), "testdata")
	testVideoPath = filepath.Join(testdataDir, "test.mp4")
}

func TestExtractFramesWholeVideo(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	framesPaths, err := ExtractFrames(tmpVideoPath, time.Duration(0.0), time.Duration(0.0))
	assert.NoError(t, err)
	assert.Greater(t, len(framesPaths), 0)
}

func TestExtractFramesWithStart(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	start := time.Duration(5 * time.Second)
	end := time.Duration(0)
	framesPaths, err := ExtractFrames(tmpVideoPath, start, end)
	assert.NoError(t, err)
	assert.Greater(t, len(framesPaths), 0)
}

func TestExtractFramesWithStartError(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	start := time.Duration(50 * time.Second) // Longer than the video
	end := time.Duration(0)
	framesPaths, err := ExtractFrames(tmpVideoPath, start, end)
	assert.Error(t, err)
	assert.Equal(t, len(framesPaths), 0)
}

func TestExtractFramesWithEnd(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	start := time.Duration(0)
	end := time.Duration(5 * time.Second)
	framesPaths, err := ExtractFrames(tmpVideoPath, start, end)
	assert.NoError(t, err)
	assert.Greater(t, len(framesPaths), 0)
}

func TestExtractFramesWithStartAndEnd(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	start := time.Duration(5 * time.Second)
	end := time.Duration(10 * time.Second)
	framesPaths, err := ExtractFrames(tmpVideoPath, start, end)
	assert.NoError(t, err)
	assert.Greater(t, len(framesPaths), 0)
}

func TestExtractFramesWithStartAndEndError(t *testing.T) {
	tmpVideoPath, tmpDir := copyVideoToTmpDir(t)
	defer cleanupTmpDir(t, tmpDir)

	start := time.Duration(10 * time.Second)
	end := time.Duration(10 * time.Second)
	framesPaths, err := ExtractFrames(tmpVideoPath, start, end)
	assert.Error(t, err)
	assert.Equal(t, len(framesPaths), 0)
}

// copyVideoToTmpDir creates a temporary directory and copies
// the test video file from the testdata folder into it.
// It returns the path to the copied video and the temporary
// directory path.
func copyVideoToTmpDir(t *testing.T) (string, string) {
	tmpDir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	source, err := os.Open(testVideoPath)
	assert.NoError(t, err)
	defer source.Close()

	tmpVideoPath := filepath.Join(tmpDir, "test.mp4")
	destination, err := os.Create(tmpVideoPath)
	assert.NoError(t, err)
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	assert.NoError(t, err)
	assert.Greater(t, nBytes, int64(0))

	return tmpVideoPath, tmpDir
}

// cleanupTmpDir removes the temporary directory created by copyVideoToTmpDir.
func cleanupTmpDir(t *testing.T, tmpDir string) {
	err := os.RemoveAll(tmpDir)
	assert.NoError(t, err)
}
