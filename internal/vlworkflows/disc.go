package vlworkflows

import "go.temporal.io/sdk/workflow"

type DiscState struct{}

func Disc(ctx workflow.Context, state *DiscState) error {
	return nil
}