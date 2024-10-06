package vltemp

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vlcontext"
	"go.temporal.io/sdk/client"
)

type VolumeWFState struct {
	Discs []string `json:"discs"`
}

type VolumeWFUpdateDiscoverNewDiscsResponse struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []string `json:"discovered"`
}

var ErrInvalidVolumeID = errors.New("invalid volume ID")

func ValidateVolumeID(volumeID string) error {
	if volumeID == "" {
		return ErrInvalidVolumeID
	}
	return nil
}

func VolumePath(ctx context.Context, volumeID string) (string, error) {
	cfg := vlcontext.GetConfig(ctx)
	if err := ValidateVolumeID(volumeID); err != nil {
		return "", err
	}
	return filepath.Join(cfg.Volume.Directory, volumeID), nil
}

func VolumeMkDir(ctx context.Context, volumeID string) error {
	dir, err := VolumePath(ctx, volumeID)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = lightOptions

func VolumeReadDiscNames(ctx context.Context, volumeID string) ([]string, error) {
	dir, err := VolumePath(ctx, volumeID)
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
	discID, err := DiscID(volumeID, discBase)
	if err != nil {
		return "", err
	}

	opts := client.StartWorkflowOptions{
		ID: discID,
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, Disc, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID:            uuid.New().String(),
		UpdateName:          DiscUpdateBootstrap,
		WorkflowID:          discID,
		WaitForStage:        client.WorkflowUpdateStageCompleted,
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
