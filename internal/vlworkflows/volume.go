package vlworkflows

import (
	"path/filepath"

	"github.com/krelinga/video-library/internal/vlactivities"
	"go.temporal.io/sdk/workflow"
)

type VolumeState struct {
	Discs []string `json:"discs"`
}

type VolumeDiscoverNewDiscsUpdateResponse struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []string
}

const VolumeDiscoverNewDiscsUpdate = "VolumeDiscoverNewDiscsUpdate"

func Volume(ctx workflow.Context, state *VolumeState) error {
	volumeName := workflow.GetInfo(ctx).WorkflowExecution.ID

	didWork := false
	if state == nil {
		// A nil state indicates that this is a freshly-created Volume,
		// so we need to initialize it and create the corresponding directory on-disk.
		state = &VolumeState{}
		err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.VolumeMkDirOptions),
			vlactivities.VolumeMkDir, volumeName).Get(ctx, nil)
		if err != nil {
			return err
		}
		didWork = true
	}

	discoverNewDiscs := func(ctx workflow.Context) (*VolumeDiscoverNewDiscsUpdateResponse, error) {
		var response VolumeDiscoverNewDiscsUpdateResponse
		var discDirs []string
		err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.VolumeReadDiscNamesOptions),
			vlactivities.VolumeReadDiscNames, volumeName).Get(ctx, &discDirs)
		if err != nil {
			return nil, err
		}
		oldDiscs := map[string]struct{}{}
		for _, disc := range state.Discs {
			oldDiscs[disc] = struct{}{}
		}
		for _, discDir := range discDirs {
			disc := filepath.Join(volumeName, discDir)
			if _, ok := oldDiscs[disc]; ok {
				continue
			}
			response.Discovered = append(response.Discovered, disc)
			state.Discs = append(state.Discs, disc)
			ctx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
				WorkflowID: disc,
			})
			err := workflow.ExecuteChildWorkflow(ctx, Disc, nil).Get(ctx, nil)
			if err != nil {
				return nil, err
			}
		}

		didWork = true
		return &response, err
	}

	err := workflow.SetUpdateHandler(ctx, VolumeDiscoverNewDiscsUpdate, discoverNewDiscs)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, func() bool {
		return didWork
	})
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Volume, state)
}
