package vlvolume

import (
	"errors"
	"testing"
	"time"

	"github.com/krelinga/video-library/internal/vlworkflows/vldisc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *WorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *WorkflowTestSuite) TestFreshlyCreated() {
	s.env.OnActivity(actConfigBased.MakeVolumeDir, mock.Anything, "test_volume").Return("/nas/media/Volumes/test_volume", nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(Workflow, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	if s.True(workflow.IsContinueAsNewError(err)) {
		var cont *workflow.ContinueAsNewError
		errors.As(err, &cont)
		conv := converter.GetDefaultDataConverter()
		actualState := &State{}
		conv.FromPayloads(cont.Input, &actualState)
		expectedState := &State{
			Directory: "/nas/media/Volumes/test_volume",
		}
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
	case *DiscoverNewDiscsResult:
		cb.a.ElementsMatch(cb.expected, s.Discovered)
	default:
		cb.a.Fail("unexpected type", success)
	}
}

func (s *WorkflowTestSuite) TestDiscoverNewDiscs() {
	s.env.RegisterWorkflow(vldisc.Workflow2)
	s.env.OnActivity(actConfigBased.ReadDiscNames, mock.Anything, "/nas/media/Volumes/test_volume").Return([]string{"disc1", "disc2"}, nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	state := &State{
		Directory: "/nas/media/Volumes/test_volume",
		Discs:     []string{"test_volume/disc1"},
	}
	s.env.RegisterDelayedCallback(func() {
		s.env.UpdateWorkflowByID("test_volume", DiscoverNewDiscs, "", &newDiscsUpdateCBs{a: s.Assertions, expected: []string{"test_volume/disc2"}}, nil)
	}, time.Hour)
	s.env.ExecuteWorkflow(Workflow, state)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	if s.True(workflow.IsContinueAsNewError(err)) {
		var cont *workflow.ContinueAsNewError
		errors.As(err, &cont)
		conv := converter.GetDefaultDataConverter()
		actualState := &State{}
		conv.FromPayloads(cont.Input, &actualState)
		expectedState := &State{
			Directory: "/nas/media/Volumes/test_volume",
			Discs:     []string{"test_volume/disc1", "test_volume/disc2"},
		}
		s.Equal(expectedState, actualState)
	}
}

func TestWorkflowTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WorkflowTestSuite))
}
