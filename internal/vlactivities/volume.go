package vlactivities

import (
	"context"
	"os"
	"strings"

	"github.com/krelinga/video-library/internal/vllib"
)

func VolumeMkDir(ctx context.Context, volumeID string) error {
	dir := vllib.VolumePath(ctx, volumeID)
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = lightOptions

func VolumeReadDiscNames(ctx context.Context, volumeID string) ([]string, error) {
	dir := vllib.VolumePath(ctx, volumeID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		out = append(out, entry.Name())
	}
	return out, nil
}

var VolumeReadDiscNamesOptions = lightOptions