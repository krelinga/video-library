package vltemp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
)

type volumeTestSuite struct {
	testSuite
}

func (s *volumeTestSuite) TestFreshlyCreated() {
	s.env.OnActivity(VolumeMkDir, mock.Anything, "test_volume").Return(nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(VolumeWF, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &VolumeWFState{})
}

func (s *volumeTestSuite) TestDiscoverNewDiscs() {
	s.env.OnActivity(VolumeReadDiscNames, mock.Anything, "test_volume").Return([]string{"disc1", "disc2"}, nil)
	s.env.OnActivity(VolumeBootstrapDisc, mock.Anything, "test_volume", "disc2").Return("test_volume/disc2", nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	state := &VolumeWFState{
		Discs: []string{"test_volume/disc1"},
	}
	s.env.RegisterDelayedCallback(
		func() {
			s.env.UpdateWorkflowByID("test_volume", VolumeWFUpdateNameDiscoverNewDiscs, "",
				assertComplete(s.Assertions, &VolumeWFUpdateDiscoverNewDiscsResponse{
					Discovered: []string{"test_volume/disc2"},
				}, nil), nil)
		}, time.Hour)
	s.env.ExecuteWorkflow(VolumeWF, state)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &VolumeWFState{
		Discs: []string{"test_volume/disc1", "test_volume/disc2"},
	})
}

func TestVolume(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(volumeTestSuite))
}
