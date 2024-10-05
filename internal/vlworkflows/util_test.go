package vlworkflows

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

func assertContinuedWithState[StateType any](a *assert.Assertions, err error, expectedState *StateType) bool {
	if !a.True(workflow.IsContinueAsNewError(err), err) {
		return false
	}
	var cont *workflow.ContinueAsNewError
	errors.As(err, &cont)
	conv := converter.GetDefaultDataConverter()
	var actualState *StateType
	conv.FromPayloads(cont.Input, &actualState)
	return a.Equal(expectedState, actualState)
}
