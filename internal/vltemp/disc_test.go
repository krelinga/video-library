package vltemp

import (
	"context"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/stretchr/testify/assert"
)

func TestDiscPath(t *testing.T) {
	tests := []struct {
		name     string
		discWfId DiscWfId
		wantPath string
		wantErr  error
	}{
		{
			name:     "valid discID",
			discWfId: DiscWfId("volume1/disc1"),
			wantPath: "/mocked/path/volume1/disc1",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{Volume: &vlconfig.Volume{Directory: "/mocked/path"}})
			gotPath, err := DiscPath(ctx, tt.discWfId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPath, gotPath)
			}
		})
	}
}
