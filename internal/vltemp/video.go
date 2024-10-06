package vltemp

import (
	"context"
	"errors"
	"path/filepath"
)

const (
	VideoUpdateBootstrap = "video-update-bootstrap"
)

type VideoLineage struct {
	FromDisc *VideoFromDisc `json:"from_disc"`
	// TODO: eventually support other options here.
}

type VideoFromDisc struct {
	DiscID   string `json:"disc_id"`
	Filename string `json:"filename"`
}

type VideoUpdateBootstrapRequest struct {
	Lineage *VideoLineage `json:"lineage"`
}

var ErrCorruptVideoLineage = errors.New("corrupt video lineage")

func VideoPath(ctx context.Context, videoLineage *VideoLineage) (string, error) {
	switch {
	case videoLineage.FromDisc != nil:
		discPath, err := DiscPath(ctx, videoLineage.FromDisc.DiscID)
		if err != nil {
			return "", errors.Join(ErrCorruptVideoLineage, err)
		}
		if videoLineage.FromDisc.Filename == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing filename"))
		}
		return filepath.Join(discPath, videoLineage.FromDisc.Filename), nil
	default:
		return "", errors.Join(ErrCorruptVideoLineage, errors.New("unknown lineage"))
	}
}

func VideoID(videoLineage *VideoLineage) (string, error) {
	switch {
	case videoLineage.FromDisc != nil:
		if videoLineage.FromDisc.DiscID == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing DiscID"))
		}
		if videoLineage.FromDisc.Filename == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing filename"))
		}
		return filepath.Join(videoLineage.FromDisc.DiscID, videoLineage.FromDisc.Filename), nil
	default:
		return "", errors.Join(ErrCorruptVideoLineage, errors.New("unknown lineage"))
	}
}
