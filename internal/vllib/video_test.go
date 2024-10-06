package vllib

import (
	"context"
	"errors"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/krelinga/video-library/internal/vltemp"
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
		videoLineage  *vltemp.VideoLineage
		expectedPath  string
		expectedError error
	}{
		{
			name: "FromDisc with valid DiscID and Filename",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					DiscID:   "volumeID/discID",
					Filename: "video.mp4",
				},
			},
			expectedPath:  "/mocked/path/volumeID/discID/video.mp4",
			expectedError: nil,
		},
		{
			name: "FromDisc with missing Filename",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					DiscID: "volumeID/discID",
				},
			},
			expectedPath:  "",
			expectedError: ErrCorruptVideoLineage,
		},
		{
			name: "FromDisc with invalid DiscID",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					DiscID:   "invalidDiscID",
					Filename: "video.mp4",
				},
			},
			expectedPath:  "",
			expectedError: ErrCorruptVideoLineage,
		},
		{
			name:          "Unknown lineage",
			videoLineage:  &vltemp.VideoLineage{},
			expectedPath:  "",
			expectedError: ErrCorruptVideoLineage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := VideoPath(ctx, tt.videoLineage)
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

func TestVideoID(t *testing.T) {
	tests := []struct {
		name          string
		videoLineage  *vltemp.VideoLineage
		expectedID    string
		expectedError error
	}{
		{
			name: "FromDisc with valid DiscID and Filename",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					DiscID:   "volumeID/discID",
					Filename: "video.mp4",
				},
			},
			expectedID:    "volumeID/discID/video.mp4",
			expectedError: nil,
		},
		{
			name: "FromDisc with missing Filename",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					DiscID: "volumeID/discID",
				},
			},
			expectedID:    "",
			expectedError: ErrCorruptVideoLineage,
		},
		{
			name: "FromDisc with missing DiscID",
			videoLineage: &vltemp.VideoLineage{
				FromDisc: &vltemp.VideoFromDisc{
					Filename: "video.mp4",
				},
			},
			expectedID:    "",
			expectedError: ErrCorruptVideoLineage,
		},
		{
			name:          "Unknown lineage",
			videoLineage:  &vltemp.VideoLineage{},
			expectedID:    "",
			expectedError: ErrCorruptVideoLineage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := VideoID(tt.videoLineage)
			if tt.expectedError != nil {
				assert.Empty(t, id)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
