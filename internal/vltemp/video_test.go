package vltemp

import (
	"context"
	"errors"
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
			name:          "FromDisc with valid DiscID and Filename",
			videoWfId:     VideoWfId("disc:volumeID/discID/video.mp4"),
			expectedPath:  "/mocked/path/volumeID/discID/video.mp4",
			expectedError: nil,
		},
		{
			name:          "FromFilepath with valid Filepath",
			videoWfId:     VideoWfId("filepath:/path/to/the/v/i/d/e/o.mp4"),
			expectedPath:  "/path/to/the/v/i/d/e/o.mp4",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := VideoPath(ctx, tt.videoWfId)
			if tt.expectedError != nil {
				assert.Empty(t, path)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}
