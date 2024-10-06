package vltemp

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

const (
	Disc                = "disc"
	DiscUpdateBootstrap = "disc-update-bootstrap"
)

type DiscState struct {
	Videos []string `json:"videos"`
}

var ErrInvalidDiscID = errors.New("invalid discID")
var ErrInvalidDiscBase = errors.New("invalid discBase")

func DiscParseID(discID string, volumeID, discBase *string) error {
	parts := strings.Split(discID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ErrInvalidDiscID
	}
	if volumeID != nil {
		*volumeID = parts[0]
	}
	if discBase != nil {
		*discBase = parts[1]
	}
	return nil
}

func DiscID(volumeID, discBase string) (string, error) {
	if err := ValidateVolumeID(volumeID); err != nil {
		return "", err
	}
	if discBase == "" {
		return "", ErrInvalidDiscBase
	}
	return filepath.Join(volumeID, discBase), nil
}

func DiscPath(ctx context.Context, discID string) (string, error) {
	var volumeID, discBase string
	err := DiscParseID(discID, &volumeID, &discBase)
	if err != nil {
		return "", err
	}
	volumePath, err := VolumePath(ctx, volumeID)
	if err != nil {
		return "", err
	}
	return filepath.Join(volumePath, discBase), nil
}
