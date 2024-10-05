package vllib

import (
	"context"
	"errors"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/krelinga/video-library/internal/vltypes"
	"github.com/stretchr/testify/assert"
)

func TestVideoPath(t *testing.T) {
	ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{
		Volume: &vlconfig.Volume{
			Directory: "/mocked/path",
		},
	})

	t.Run("FromDisc with valid DiscID and Filename", func(t *testing.T) {
		videoLineage := &vltypes.VideoLineage{
			FromDisc: &vltypes.VideoFromDisc{
				DiscID:   "volumeID/discID",
				Filename: "video.mp4",
			},
		}

		expectedPath := "/mocked/path/volumeID/discID/video.mp4"

		path, err := VideoPath(ctx, videoLineage)
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
	})

	t.Run("FromDisc with missing Filename", func(t *testing.T) {
		videoLineage := &vltypes.VideoLineage{
			FromDisc: &vltypes.VideoFromDisc{
				DiscID: "volumeID/discID",
			},
		}

		path, err := VideoPath(ctx, videoLineage)
		assert.Empty(t, path)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrCorruptVideoLineage))
	})

	t.Run("FromDisc with invalid DiscID", func(t *testing.T) {
		videoLineage := &vltypes.VideoLineage{
			FromDisc: &vltypes.VideoFromDisc{
				DiscID:   "invalidDiscID",
				Filename: "video.mp4",
			},
		}

		path, err := VideoPath(ctx, videoLineage)
		assert.Empty(t, path)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrCorruptVideoLineage))
	})

	t.Run("Unknown lineage", func(t *testing.T) {
		videoLineage := &vltypes.VideoLineage{}

		path, err := VideoPath(ctx, videoLineage)
		assert.Empty(t, path)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrCorruptVideoLineage))
	})
}
