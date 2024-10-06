package vlworkflows

import (
	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vltemp"
	"go.temporal.io/sdk/workflow"
)

func DiscWF(ctx workflow.Context, state *vltemp.DiscWFState) error {
	discId := workflow.GetInfo(ctx).WorkflowExecution.ID
	wt := workTracker{}

	bootstrap := func(ctx workflow.Context) (err error) {
		defer wt.WorkIfNoError(err)

		state = &vltemp.DiscWFState{}
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

	return workflow.NewContinueAsNewError(ctx, DiscWF, state)
}
