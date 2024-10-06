package vlworkflows

import (
	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vltemp"
	"github.com/krelinga/video-library/internal/vltypes"
	"go.temporal.io/sdk/workflow"
)

func Disc(ctx workflow.Context, state *vltypes.DiscState) error {
	discId := workflow.GetInfo(ctx).WorkflowExecution.ID
	wt := workTracker{}

	bootstrap := func(ctx workflow.Context) (err error) {
		defer wt.WorkIfNoError(err)

		state = &vltypes.DiscState{}
		var videoFiles []string
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.DiscReadVideoNamesOptions),
			vlactivities.DiscReadVideoNames, discId).Get(ctx, &videoFiles)
		if err != nil {
			return
		}
		for _, videoFile := range state.Videos {
			var videoId string
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, vlactivities.DiscBootstrapVideoOptions),
				vlactivities.DiscBootstrapVideo, discId, videoFile).Get(ctx, &videoId)
			if err != nil {
				return err
			}
			state.Videos = append(state.Videos, videoId)
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, vltemp.DiscUpdateBootstrap, bootstrap)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Disc, state)
}
