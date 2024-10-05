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

type updateCallbacks[outType any] struct {
	accept func()
	reject func(error)
	complete func(*outType, error)
}

func (cb *updateCallbacks[outType]) Accept() {
	if cb.accept != nil {
		cb.accept()
	}
}

func (cb *updateCallbacks[outType]) Reject(err error) {
	if cb.reject != nil {
		cb.reject(err)
	}
}

func (cb *updateCallbacks[outType]) Complete(success any, err error) {
	if cb.complete != nil {
		cb.complete(success.(*outType), err)
	}
}

func assertComplete[outType any](a *assert.Assertions, out *outType, err error) *updateCallbacks[outType] {
	return &updateCallbacks[outType]{
		complete: func(actual *outType, actualErr error) {
			a.Equal(out, actual)
			a.ErrorIs(err, actualErr)
		},
	}
}