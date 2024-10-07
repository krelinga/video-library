package vltemp

import (
	"context"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/stretchr/testify/assert"
)

func TestDiscParseID(t *testing.T) {
	tests := []struct {
		name         string
		discID       string
		volumeID     *string
		discBase     *string
		wantErr      error
		wantVolumeID string
		wantDiscBase string
	}{
		{
			name:         "valid discID with non-nil volumeID and discBase",
			discID:       "volume1/disc1",
			volumeID:     new(string),
			discBase:     new(string),
			wantErr:      nil,
			wantVolumeID: "volume1",
			wantDiscBase: "disc1",
		},
		{
			name:         "valid discID with nil volumeID",
			discID:       "volume2/disc2",
			volumeID:     nil,
			discBase:     new(string),
			wantErr:      nil,
			wantDiscBase: "disc2",
		},
		{
			name:         "valid discID with nil discBase",
			discID:       "volume3/disc3",
			volumeID:     new(string),
			discBase:     nil,
			wantErr:      nil,
			wantVolumeID: "volume3",
		},
		{
			name:     "invalid discID with missing slash",
			discID:   "volume4disc4",
			volumeID: new(string),
			discBase: new(string),
			wantErr:  ErrInvalidDiscID,
		},
		{
			name:     "invalid discID with empty volume",
			discID:   "/disc5",
			volumeID: new(string),
			discBase: new(string),
			wantErr:  ErrInvalidDiscID,
		},
		{
			name:     "invalid discID with empty disc",
			discID:   "volume6/",
			volumeID: new(string),
			discBase: new(string),
			wantErr:  ErrInvalidDiscID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DiscParseID(tt.discID, tt.volumeID, tt.discBase)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				if tt.volumeID != nil {
					assert.Equal(t, tt.wantVolumeID, *tt.volumeID)
				}
				if tt.discBase != nil {
					assert.Equal(t, tt.wantDiscBase, *tt.discBase)
				}
			}
		})
	}
}

func TestDiscPathWithMockedDiscParseID(t *testing.T) {
	tests := []struct {
		name     string
		discID   string
		wantPath string
		wantErr  error
	}{
		{
			name:     "valid discID",
			discID:   "volume1/disc1",
			wantPath: "/mocked/path/volume1/disc1",
			wantErr:  nil,
		},
		{
			name:     "invalid discID",
			discID:   "volume2disc2",
			wantPath: "",
			wantErr:  ErrInvalidDiscID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{Volume: &vlconfig.Volume{Directory: "/mocked/path"}})
			gotPath, err := DiscPath(ctx, tt.discID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPath, gotPath)
			}
		})
	}
}
