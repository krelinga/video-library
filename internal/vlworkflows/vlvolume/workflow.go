package vlvolume

import (
	"time"

	"github.com/krelinga/video-library/internal/vlactivities"
	"go.temporal.io/sdk/workflow"
)

var actConfigBased *vlactivities.ConfigBased = nil

func Workflow(ctx workflow.Context, state *State) error {
	volumeName := workflow.GetInfo(ctx).WorkflowExecution.ID
	// TODO: Setting this globally does not make sense to me.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})
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
	}

	return workflow.NewContinueAsNewError(ctx, Workflow, state)
}
