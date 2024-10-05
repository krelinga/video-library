package vllib

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

var ErrInvalidDiscID = errors.New("invalid disc name")

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

func DiscPath(ctx context.Context, discID string) (string, error) {
	var volumeID, discBase string
	err := DiscParseID(discID, &volumeID, &discBase)
	if err != nil {
		return "", err
	}
	return filepath.Join(VolumePath(ctx, volumeID), discBase), nil
}