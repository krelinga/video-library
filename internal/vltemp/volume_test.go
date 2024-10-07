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
	s.env.OnActivity(VolumeMkDir, mock.Anything, VolumeWfId("test_volume")).Return(nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(VolumeWF, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &VolumeWFState{})
}

func (s *volumeTestSuite) TestDiscoverNewDiscs() {
	s.env.OnActivity(VolumeReadDiscNames, mock.Anything, VolumeWfId("test_volume")).Return([]string{"disc1", "disc2"}, nil)
	s.env.OnActivity(VolumeBootstrapDisc, mock.Anything, VolumeWfId("test_volume"), "disc2").Return(DiscWfId("test_volume/disc2"), nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	state := &VolumeWFState{
		Discs: []DiscWfId{DiscWfId("test_volume/disc1")},
	}
	s.env.RegisterDelayedCallback(
		func() {
			s.env.UpdateWorkflowByID("test_volume", VolumeWFUpdateNameDiscoverNewDiscs, "",
				assertComplete(s.Assertions, &VolumeWFUpdateDiscoverNewDiscsResponse{
					Discovered: []DiscWfId{DiscWfId("test_volume/disc2")},
				}, nil), nil)
		}, time.Hour)
	s.env.ExecuteWorkflow(VolumeWF, state)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &VolumeWFState{
		Discs: []DiscWfId{DiscWfId("test_volume/disc1"), DiscWfId("test_volume/disc2")},
	})
}

func TestVolumeWf(t *testing.T) {
	suite.Run(t, new(volumeTestSuite))
}
