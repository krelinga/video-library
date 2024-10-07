package vltemp

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vlcontext"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type VolumeWFState struct {
	Discs []DiscWfId `json:"discs"`
}

type VolumeWFUpdateDiscoverNewDiscsResponse struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []DiscWfId `json:"discovered"`
}

func VolumePath(ctx context.Context, volumeWfId VolumeWfId) string {
	cfg := vlcontext.GetConfig(ctx)
	return filepath.Join(cfg.Volume.Directory, volumeWfId.Name())
}

func VolumeMkDir(ctx context.Context, volumeWfId VolumeWfId) error {
	dir := VolumePath(ctx, volumeWfId)
	return os.MkdirAll(dir, 0755)
}

var VolumeMkDirOptions = lightOptions

func VolumeReadDiscNames(ctx context.Context, volumeWfId VolumeWfId) ([]string, error) {
	dir := VolumePath(ctx, volumeWfId)
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

func VolumeBootstrapDisc(ctx context.Context, volumeWfId VolumeWfId, discFilename string) (DiscWfId, error) {
	temporalClient := vlcontext.GetTemporalClient(ctx)
	discWfId, err := NewDiscWfId(volumeWfId, discFilename)
	if err != nil {
		return "", err
	}

	opts := client.StartWorkflowOptions{
		ID: string(discWfId),
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, DiscWF, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID:            uuid.New().String(),
		UpdateName:          DiscWFUpdateNameBootstrap,
		WorkflowID:          string(discWfId),
		WaitForStage:        client.WorkflowUpdateStageCompleted,
		FirstExecutionRunID: wf.GetRunID(),
	})
	if err != nil {
		return "", err
	}
	if err := updateHandle.Get(ctx, nil); err != nil {
		return "", err
	}

	return discWfId, nil
}

var VolumeBootstrapDiscOptions = lightOptions

const VolumeWFUpdateNameDiscoverNewDiscs = "VolumeWFUpdateDiscoverNewDiscs"

func VolumeWF(ctx workflow.Context, state *VolumeWFState) error {
	volumeWfId := VolumeWfId(workflow.GetInfo(ctx).WorkflowExecution.ID)
	if err := volumeWfId.Validate(); err != nil {
		return err
	}

	wt := workTracker{}
	if state == nil {
		// A nil state indicates that this is a freshly-created Volume,
		// so we need to initialize it and create the corresponding directory on-disk.
		state = &VolumeWFState{}
		err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, VolumeMkDirOptions),
			VolumeMkDir, volumeWfId).Get(ctx, nil)
		if err != nil {
			return err
		}
		wt.Work()
	}

	discoverNewDiscs := func(ctx workflow.Context) (response *VolumeWFUpdateDiscoverNewDiscsResponse, err error) {
		defer wt.WorkIfNoError(err)
		var discDirs []string
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, VolumeReadDiscNamesOptions),
			VolumeReadDiscNames, volumeWfId).Get(ctx, &discDirs)
		if err != nil {
			return
		}
		oldDiscs := map[DiscWfId]struct{}{}
		for _, disc := range state.Discs {
			oldDiscs[disc] = struct{}{}
		}
		for _, discDir := range discDirs {
			var discWfId DiscWfId
			discWfId, err = NewDiscWfId(volumeWfId, discDir)
			if err != nil {
				return
			}
			if _, ok := oldDiscs[discWfId]; ok {
				continue
			}
			if response == nil {
				response = &VolumeWFUpdateDiscoverNewDiscsResponse{}
			}
			response.Discovered = append(response.Discovered, discWfId)
			state.Discs = append(state.Discs, discWfId)
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, VolumeBootstrapDiscOptions),
				VolumeBootstrapDisc, volumeWfId, discDir).Get(ctx, nil)
			if err != nil {
				return
			}
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, VolumeWFUpdateNameDiscoverNewDiscs, discoverNewDiscs)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, VolumeWF, state)
}

// A more-refined string to handle Temporal Workflow IDs for Volume workflows.
//
// Use NewVolumeWfId() to create a new VolumeWfId.  You can also directly case from a string
// with `VolumeWfId("my-volume")`, but this will not validate the ID.  You can validate the ID
// with the Validate() method.  Any other methods called on an invalid VolumeWfId will panic.
type VolumeWfId string

// Checks if the VolumeWfId is valid.
func (id VolumeWfId) Validate() error {
	if !nameIsValid(string(id)) {
		return ErrInvalidWorkflowId
	}
	return nil
}

// Returns the name of the Volume.
//
// Panics if the VolumeWfId is invalid.
func (id VolumeWfId) Name() string {
	if err := id.Validate(); err != nil {
		panic(err)
	}
	return string(id)
}

// Returns a VolumeWfId for the given workflow name, or an error if the name is invalid.
func NewVolumeWfId(name string) (VolumeWfId, error) {
	id := VolumeWfId(name)
	return id, id.Validate()
}
