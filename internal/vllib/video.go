package vllib

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/krelinga/video-library/internal/vltypes"
)

var ErrCorruptVideoLineage = errors.New("corrupt video lineage")

func VideoPath(ctx context.Context, videoLineage *vltypes.VideoLineage) (string, error) {
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

func VideoID(videoLineage *vltypes.VideoLineage) (string, error) {
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
