package vlactivities

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/krelinga/video-library/internal/vllib"
	"github.com/krelinga/video-library/internal/vlqueues"
	"go.temporal.io/sdk/workflow"
)

func VolumeMkDir(ctx context.Context, volumeID string) error {
	dir := vllib.VolumePath(ctx, volumeID)
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           vlqueues.Light,
}

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

var VolumeReadDiscNamesOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           vlqueues.Light,
}
