package vllib

import (
	"errors"
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
