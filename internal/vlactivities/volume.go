package vlactivities

import (
	"context"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/krelinga/video-library/internal/vllib"
	"github.com/krelinga/video-library/internal/vlqueues"
	"go.temporal.io/sdk/client"
)

func VolumeMkDir(ctx context.Context, volumeID string) error {
	dir, err := vllib.VolumePath(ctx, volumeID)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = lightOptions

func VolumeReadDiscNames(ctx context.Context, volumeID string) ([]string, error) {
	dir, err := vllib.VolumePath(ctx, volumeID)
	if err != nil {
		return nil, err
	}
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

func VolumeBootstrapDisc(ctx context.Context, volumeID, discBase string) (string, error) {
	temporalClient := vlcontext.GetTemporalClient(ctx)
	discID, err := vllib.DiscID(volumeID, discBase)
	if err != nil {
		return "", err
	}

	opts := client.StartWorkflowOptions{
		ID: discID,
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, vlqueues.Disc, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID: uuid.New().String(),
		UpdateName: vlqueues.DiscUpdateBootstrap,
		WorkflowID: discID,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		FirstExecutionRunID: wf.GetRunID(),
	})
	if err != nil {
		return "", err
	}
	if err := updateHandle.Get(ctx, nil); err != nil {
		return "", err
	}

	return discID, nil
}

var VolumeBootstrapDiscOptions = lightOptions
