package vlworkflows

import (
	"path/filepath"

	"github.com/krelinga/video-library/internal/vltemp"
	"go.temporal.io/sdk/workflow"
)

const VolumeWFUpdateDiscoverNewDiscsName = "VolumeWFUpdateDiscoverNewDiscs"

func VolumeWF(ctx workflow.Context, state *vltemp.VolumeWFState) error {
	volumeID := workflow.GetInfo(ctx).WorkflowExecution.ID

	wt := workTracker{}
	if state == nil {
		// A nil state indicates that this is a freshly-created Volume,
		// so we need to initialize it and create the corresponding directory on-disk.
		state = &vltemp.VolumeWFState{}
		err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vltemp.VolumeMkDirOptions),
			vltemp.VolumeMkDir, volumeID).Get(ctx, nil)
		if err != nil {
			return err
		}
		wt.Work()
	}

	discoverNewDiscs := func(ctx workflow.Context) (response *vltemp.VolumeWFUpdateDiscoverNewDiscsResponse, err error) {
		defer wt.WorkIfNoError(err)
		var discDirs []string
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vltemp.VolumeReadDiscNamesOptions),
			vltemp.VolumeReadDiscNames, volumeID).Get(ctx, &discDirs)
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
				response = &vltemp.VolumeWFUpdateDiscoverNewDiscsResponse{}
			}
			response.Discovered = append(response.Discovered, disc)
			state.Discs = append(state.Discs, disc)
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, vltemp.VolumeBootstrapDiscOptions),
				vltemp.VolumeBootstrapDisc, volumeID, discDir).Get(ctx, nil)
			if err != nil {
				return
			}
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, VolumeWFUpdateDiscoverNewDiscsName, discoverNewDiscs)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, VolumeWF, state)
}
