package vltemp

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/krelinga/video-library/internal/vlcontext"
)

type VolumeState struct {
	Discs []string `json:"discs"`
}

type VolumeDiscoverNewDiscsUpdateResponse struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []string `json:"discovered"`
}

var ErrInvalidVolumeID = errors.New("invalid volume ID")

func ValidateVolumeID(volumeID string) error {
	if volumeID == "" {
		return ErrInvalidVolumeID
	}
	return nil
}

func VolumePath(ctx context.Context, volumeID string) (string, error) {
	cfg := vlcontext.GetConfig(ctx)
	if err := ValidateVolumeID(volumeID); err != nil {
		return "", err
	}
	return filepath.Join(cfg.Volume.Directory, volumeID), nil
}