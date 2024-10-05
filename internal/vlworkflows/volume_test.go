package vlworkflows

import (
	"errors"
	"testing"
	"time"

	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

type VolumeTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *VolumeTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *VolumeTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *VolumeTestSuite) TestFreshlyCreated() {
	s.env.OnActivity(vlactivities.VolumeMkDir, mock.Anything, "test_volume").Return(nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(Volume, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	if s.True(workflow.IsContinueAsNewError(err)) {
		var cont *workflow.ContinueAsNewError
		errors.As(err, &cont)
		conv := converter.GetDefaultDataConverter()
		actualState := &VolumeState{}
		conv.FromPayloads(cont.Input, &actualState)
		expectedState := &VolumeState{}
		s.Equal(expectedState, actualState)
	}
}

type newDiscsUpdateCBs struct {
	a        *assert.Assertions
	expected []string
}

func (cb *newDiscsUpdateCBs) Accept() {}

func (cb *newDiscsUpdateCBs) Reject(err error) {
	cb.a.NoError(err)
}

func (cb *newDiscsUpdateCBs) Complete(success any, err error) {
	if !cb.a.NoError(err) {
		return
	}
	switch s := success.(type) {
	case *VolumeDiscoverNewDiscsUpdateResponse:
		cb.a.ElementsMatch(cb.expected, s.Discovered)
	default:
		cb.a.Fail("unexpected type", success)
	}
}

func (s *VolumeTestSuite) TestDiscoverNewDiscs() {
	s.env.RegisterWorkflow(Disc)
	s.env.OnActivity(vlactivities.VolumeReadDiscNames, mock.Anything, "test_volume").Return([]string{"disc1", "disc2"}, nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	state := &VolumeState{
		Discs: []string{"test_volume/disc1"},
	}
	s.env.RegisterDelayedCallback(
		func() {
			s.env.UpdateWorkflowByID("test_volume", VolumeDiscoverNewDiscsUpdate, "",
				&newDiscsUpdateCBs{a: s.Assertions, expected: []string{"test_volume/disc2"}}, nil)
		}, time.Hour)
	s.env.ExecuteWorkflow(Volume, state)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	if s.True(workflow.IsContinueAsNewError(err)) {
		var cont *workflow.ContinueAsNewError
		errors.As(err, &cont)
		conv := converter.GetDefaultDataConverter()
		actualState := &VolumeState{}
		conv.FromPayloads(cont.Input, &actualState)
		expectedState := &VolumeState{
			Discs: []string{"test_volume/disc1", "test_volume/disc2"},
		}
		s.Equal(expectedState, actualState)
	}
}

func TestWorkflowTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VolumeTestSuite))
}
