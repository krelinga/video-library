package vlworkflows

import (
	"github.com/krelinga/video-library/internal/vlactivities"
	"go.temporal.io/sdk/workflow"
)

type DiscState struct {
	Videos []string `json:"videos"`
}

const DiscBootstrapUpdate = "DiscBootstrapUpdate"

func Disc(ctx workflow.Context, state *DiscState) error {
	discId := workflow.GetInfo(ctx).WorkflowExecution.ID
	wt := workTracker{}

	bootstrap := func(ctx workflow.Context) (err error) {
		defer wt.WorkIfNoError(err)

		state = &DiscState{}
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, vlactivities.DiscReadVideoNamesOptions),
			vlactivities.DiscReadVideoNames, discId).Get(ctx, &state.Videos)
		if err != nil {
			return
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, DiscBootstrapUpdate, bootstrap)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, Disc, state)
}
