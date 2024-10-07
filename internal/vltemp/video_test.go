package vltemp

import (
	"context"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/stretchr/testify/assert"
)

func TestVideoPath(t *testing.T) {
	ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{
		Volume: &vlconfig.Volume{
			Directory: "/mocked/path",
		},
	})

	tests := []struct {
		name          string
		videoWfId     VideoWfId
		expectedPath  string
		expectedError error
	}{
		{
			name:         "FromDisc with valid DiscID and Filename",
			videoWfId:    VideoWfId("disc:volumeID/discID/video.mp4"),
			expectedPath: "/mocked/path/volumeID/discID/video.mp4",
		},
		{
			name:         "FromFilepath with valid Filepath",
			videoWfId:    VideoWfId("filepath:/path/to/the/v/i/d/e/o.mp4"),
			expectedPath: "/path/to/the/v/i/d/e/o.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := VideoPath(ctx, tt.videoWfId)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}
func TestNewVideoWfIdFromDisc(t *testing.T) {
	tests := []struct {
		name          string
		discWfId      DiscWfId
		videoFileName string
		expectedId    VideoWfId
		expectedError error
	}{
		{
			name:          "Valid DiscWfId and VideoFileName",
			discWfId:      DiscWfId("volumeID/discID"),
			videoFileName: "video.mp4",
			expectedId:    VideoWfId("disc:volumeID/discID/video.mp4"),
			expectedError: nil,
		},
		{
			name:          "Invalid VideoFileName with subdirectory",
			discWfId:      DiscWfId("volumeID/discID"),
			videoFileName: "invalid/video.mp4",
			expectedId:    "",
			expectedError: ErrInvalidWorkflowId,
		},
		{
			name:          "Invalid VideoFileName with colon",
			discWfId:      DiscWfId("volumeID/discID"),
			videoFileName: "invalid:video.mp4",
			expectedId:    "",
			expectedError: ErrInvalidWorkflowId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewVideoWfIdFromDisc(tt.discWfId, tt.videoFileName)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Panics(t, func() { id.Protocol() })
				assert.Panics(t, func() { id.FromDisc() })
				assert.Panics(t, func() { id.FromFilepath() })
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, id)
			}
		})
	}
}
func TestNewVideoWfIdFromFilepath(t *testing.T) {
	tests := []struct {
		name          string
		filepath      string
		expectedId    VideoWfId
		expectedError error
	}{
		{
			name:          "Valid Filepath",
			filepath:      "/path/to/video.mp4",
			expectedId:    VideoWfId("filepath:/path/to/video.mp4"),
			expectedError: nil,
		},
		{
			name:          "Invalid Filepath with colon",
			filepath:      "/path/to/invalid:video.mp4",
			expectedId:    "",
			expectedError: ErrInvalidWorkflowId,
		},
		{
			name:          "Invalid Filepath with relative path",
			filepath:      "relative/path/to/video.mp4",
			expectedId:    "",
			expectedError: ErrInvalidWorkflowId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewVideoWfIdFromFilepath(tt.filepath)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Panics(t, func() { id.Protocol() })
				assert.Panics(t, func() { id.FromDisc() })
				assert.Panics(t, func() { id.FromFilepath() })
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, id)
			}
		})
	}
}
