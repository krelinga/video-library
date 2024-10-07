package vltemp

import (
	"context"
	"fmt"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/stretchr/testify/assert"
)

func TestDiscPath(t *testing.T) {
	ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{Volume: &vlconfig.Volume{Directory: "/mocked/path"}})
	gotPath := DiscPath(ctx, DiscWfId("volume1/disc1"))
	assert.Equal(t, "/mocked/path/volume1/disc1", gotPath)
}
func TestNewDiscWfId(t *testing.T) {
	tests := []struct {
		volumeWfId   VolumeWfId
		discFilename string
		expectedId   DiscWfId
		expectError  bool
	}{
		{
			volumeWfId:   VolumeWfId("volume1"),
			discFilename: "disc1",
			expectedId:   DiscWfId("volume1/disc1"),
			expectError:  false,
		},
		{
			volumeWfId:   VolumeWfId("volume2"),
			discFilename: "disc2",
			expectedId:   DiscWfId("volume2/disc2"),
			expectError:  false,
		},
		{
			volumeWfId:   VolumeWfId("volume3"),
			discFilename: "",
			expectError:  true,
		},
		{
			volumeWfId:   "",
			discFilename: "disc4",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.volumeWfId, tt.discFilename), func(t *testing.T) {
			id, err := NewDiscWfId(tt.volumeWfId, tt.discFilename)
			if tt.expectError {
				assert.ErrorIs(t, err, ErrInvalidWorkflowId)
				assert.Panics(t, func() { id.Name() })
				assert.Panics(t, func() { id.VolumeWfId() })
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, id)
			}
		})
	}
}
