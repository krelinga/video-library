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

func VolumeMkDir(ctx context.Context, volumeName string) error {
	dir := vllib.VolumePath(ctx, volumeName)
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           vlqueues.Light,
}

func VolumeReadDiscNames(ctx context.Context, volumeName string) ([]string, error) {
	dir := vllib.VolumePath(ctx, volumeName)
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
