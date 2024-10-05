package vlactivities

import (
	"context"
	"os"
	"path/filepath"

	"github.com/krelinga/video-library/internal/vlconfig"
)

type ConfigBased struct {
	config *vlconfig.Root
}

func (cb *ConfigBased) MakeVolumeDir(ctx context.Context, workflowName string) (string, error) {
	dir := filepath.Join(cb.config.Volume.Directory, workflowName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func (cb *ConfigBased) ReadDiscNames(ctx context.Context, volumeDir string) ([]string, error) {
	entries, err := os.ReadDir(volumeDir)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		out = append(out, entry.Name())
	}
	return out, nil
}
