package vllib

import (
	"context"
	"path/filepath"

	"github.com/krelinga/video-library/internal/vlcontext"
)

func VolumePath(ctx context.Context, volumeID string) string {
	cfg := vlcontext.GetConfig(ctx)
	return filepath.Join(cfg.Volume.Directory, volumeID)
}