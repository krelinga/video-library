package ids

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiscWfId(t *testing.T) {
	tests := []struct {
		name             string
		volumeWfIdString string
		discName         string
		wantErr          error
	}{
		{
			name:             "valid discWfId",
			volumeWfIdString: "volume1",
			discName:         "disc1",
			wantErr:          nil,
		},
		{
			name:             "invalid discWfId with empty discName",
			volumeWfIdString: "volume2",
			discName:         "",
			wantErr:          ErrInvalidDiscNameString,
		},
		{
			name:             "invalid discWfId with slash in discName",
			volumeWfIdString: "volume3",
			discName:         "disc/3",
			wantErr:          ErrInvalidDiscNameString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volumeWfId, err := NewVolumeWfId(tt.volumeWfIdString)
			if !assert.NoError(t, err) {
				return
			}
			got, err := NewDiscWfId(volumeWfId, tt.discName)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				if got != nil {
					assert.Equal(t, volumeWfId, got.VolumeWfId())
					assert.Equal(t, tt.discName, got.DiscName())
				}
			}
		})
	}
}

func TestNewDiscWfIdFromString(t *testing.T) {
	tests := []struct {
		name                 string
		asString             string
		wantErr              error
		wantVolumeWFIdString string
		wantDiscName         string
	}{
		{
			name:                 "valid discWfId string",
			asString:             "volume1/disc1",
			wantErr:              nil,
			wantVolumeWFIdString: "volume1",
			wantDiscName:         "disc1",
		},
		{
			name:     "invalid discWfId string with missing slash",
			asString: "volume2disc2",
			wantErr:  ErrInvalidDiscWfIdString,
		},
		{
			name:     "invalid discWfId string with empty parts",
			asString: "volume3/",
			wantErr:  ErrInvalidDiscNameString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscWfIdFromString(tt.asString)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.asString, got.String())
				wantVolumeWfId, err := NewVolumeWfId(tt.wantVolumeWFIdString)
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, wantVolumeWfId, got.VolumeWfId())
				assert.Equal(t, tt.wantDiscName, got.DiscName())
			}
		})
	}
}
