package vldisc

import "go.temporal.io/sdk/workflow"

// TODO: this is named Workflow2 because otherwise we get a conflict with the Workflow function in vlvolume.
func Workflow2(ctx workflow.Context, state *State) error {
	return nil
}
