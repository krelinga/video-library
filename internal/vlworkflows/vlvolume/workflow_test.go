package vlvolume

import (
	"errors"
	"testing"

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

func (s *WorkflowTestSuite) TestHappyPath() {
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

func TestWorkflowTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WorkflowTestSuite))
}
