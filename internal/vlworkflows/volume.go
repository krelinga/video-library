package vlworkflows

import (
	"path/filepath"

	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vltypes"
	"go.temporal.io/sdk/workflow"
)

const VolumeDiscoverNewDiscsUpdate = "VolumeDiscoverNewDiscsUpdate"

func Volume(ctx workflow.Context, state *vltypes.VolumeState) error {
	volumeID := workflow.GetInfo(ctx).WorkflowExecution.ID

	wt := workTracker{}
	if state == nil {
		// A nil state indicates that this is a freshly-created Volume,
		// so we need to initialize it and create the corresponding directory on-disk.
		state = &vltypes.VolumeState{}
		err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.VolumeMkDirOptions),
			vlactivities.VolumeMkDir, volumeID).Get(ctx, nil)
		if err != nil {
			return err
		}
		wt.Work()
	}

	discoverNewDiscs := func(ctx workflow.Context) (response *vltypes.VolumeDiscoverNewDiscsUpdateResponse, err error) {
		defer wt.WorkIfNoError(err)
		var discDirs []string
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.VolumeReadDiscNamesOptions),
			vlactivities.VolumeReadDiscNames, volumeID).Get(ctx, &discDirs)
		if err != nil {
			return
		}
		oldDiscs := map[string]struct{}{}
		for _, disc := range state.Discs {
			oldDiscs[disc] = struct{}{}
		}
		for _, discDir := range discDirs {
			disc := filepath.Join(volumeID, discDir)
			if _, ok := oldDiscs[disc]; ok {
				continue
			}
			if response == nil {
				response = &vltypes.VolumeDiscoverNewDiscsUpdateResponse{}
			}
			response.Discovered = append(response.Discovered, disc)
			state.Discs = append(state.Discs, disc)
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, vlactivities.VolumeBootstrapDiscOptions),
				vlactivities.VolumeBootstrapDisc, volumeID, discDir).Get(ctx, nil)
			if err != nil {
				return
			}
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, VolumeDiscoverNewDiscsUpdate, discoverNewDiscs)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Volume, state)
}
