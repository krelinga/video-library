package vlworkflows

import (
	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vlconst"
	"github.com/krelinga/video-library/internal/vltypes"
	"go.temporal.io/sdk/workflow"
)

func Disc(ctx workflow.Context, state *vltypes.DiscState) error {
	discId := workflow.GetInfo(ctx).WorkflowExecution.ID
	wt := workTracker{}

	bootstrap := func(ctx workflow.Context) (err error) {
		defer wt.WorkIfNoError(err)

		state = &vltypes.DiscState{}
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.DiscReadVideoNamesOptions),
			vlactivities.DiscReadVideoNames, discId).Get(ctx, &state.Videos)
		if err != nil {
			return
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, vlconst.DiscUpdateBootstrap, bootstrap)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Disc, state)
}
