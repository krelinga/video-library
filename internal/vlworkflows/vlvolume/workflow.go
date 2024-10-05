package vlvolume

import (
	"path/filepath"
	"time"

	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vlworkflows/vldisc"
	"go.temporal.io/sdk/workflow"
)

var actConfigBased *vlactivities.ConfigBased = nil

type DiscoverNewDiscsResult struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []string
}

const DiscoverNewDiscs = "DiscoverNewDiscs"

func Workflow(ctx workflow.Context, state *State) error {
	volumeName := workflow.GetInfo(ctx).WorkflowExecution.ID
	// TODO: Setting this globally does not make sense to me.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	didWork := false
	if state == nil {
		// A nil state indicates that this is a freshly-created Volume,
		// so we need to initialize it and create the corresponding directory on-disk.
		state = &State{}
		var dir string
		err := workflow.ExecuteActivity(ctx, actConfigBased.MakeVolumeDir, volumeName).Get(ctx, &dir)
		if err != nil {
			return err
		}
		state.Directory = dir
		didWork = true
	}

	discoverNewDiscs := func(ctx workflow.Context) (*DiscoverNewDiscsResult, error) {
		var result DiscoverNewDiscsResult
		var discDirs []string
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		})
		err := workflow.ExecuteActivity(ctx, actConfigBased.ReadDiscNames, state.Directory).Get(ctx, &discDirs)
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
			result.Discovered = append(result.Discovered, disc)
			state.Discs = append(state.Discs, disc)
			ctx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
				WorkflowID: disc,
			})
			err := workflow.ExecuteChildWorkflow(ctx, vldisc.Workflow2, nil).Get(ctx, nil)
			if err != nil {
				return nil, err
			}
		}

		didWork = true
		return &result, err
	}

	err := workflow.SetUpdateHandler(ctx, DiscoverNewDiscs, discoverNewDiscs)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, func() bool {
		return didWork
	})
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Workflow, state)
}
