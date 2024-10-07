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
func TestNewDiscWFID(t *testing.T) {
	tests := []struct {
		name        string
		volumeWFID  VolumeWFID
		discDirName string
		wantErr     error
		wantValid   bool
	}{
		{
			name:        "valid discWFID",
			volumeWFID:  VolumeWFID("volume1"),
			discDirName: "disc1",
			wantErr:     nil,
		},
		{
			name:        "invalid discWFID with empty discDirName",
			volumeWFID:  VolumeWFID("volume2"),
			discDirName: "",
			wantErr:     ErrInvalidDiscDirName,
		},
		{
			name:        "invalid discWFID with slash in discDirName",
			volumeWFID:  VolumeWFID("volume3"),
			discDirName: "disc/3",
			wantErr:     ErrInvalidDiscDirName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscWFID(tt.volumeWFID, tt.discDirName)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.True(t, got.valid)
				assert.Equal(t, tt.volumeWFID, got.VolumeWFID())
				assert.Equal(t, tt.discDirName, got.DiscDirName())
			}
		})
	}
}

func TestNewDiscWFIDFromString(t *testing.T) {
	tests := []struct {
		name      string
		asString  string
		wantErr   error
		wantValid bool
	}{
		{
			name:      "valid discWFID string",
			asString:  "volume1/disc1",
			wantErr:   nil,
			wantValid: true,
		},
		{
			name:      "invalid discWFID string with missing slash",
			asString:  "volume2disc2",
			wantErr:   ErrInvalidDiscWFIDString,
			wantValid: false,
		},
		{
			name:      "invalid discWFID string with empty parts",
			asString:  "volume3/",
			wantErr:   ErrInvalidDiscDirName,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscWFIDFromString(tt.asString)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, got.valid)
			}
		})
	}
}
