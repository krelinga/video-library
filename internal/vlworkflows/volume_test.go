package vlworkflows

import (
	"testing"
	"time"

	"github.com/krelinga/video-library/internal/vlactivities"
	"github.com/krelinga/video-library/internal/vltypes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
)

type volumeTestSuite struct {
	testSuite
}

func (s *volumeTestSuite) TestFreshlyCreated() {
	s.env.OnActivity(vlactivities.VolumeMkDir, mock.Anything, "test_volume").Return(nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(Volume, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &vltypes.VolumeState{})
}

func (s *volumeTestSuite) TestDiscoverNewDiscs() {
	s.env.OnActivity(vlactivities.VolumeReadDiscNames, mock.Anything, "test_volume").Return([]string{"disc1", "disc2"}, nil)
	s.env.OnActivity(vlactivities.VolumeBootstrapDisc, mock.Anything, "test_volume", "disc2").Return("test_volume/disc2", nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	state := &vltypes.VolumeState{
		Discs: []string{"test_volume/disc1"},
	}
	s.env.RegisterDelayedCallback(
		func() {
			s.env.UpdateWorkflowByID("test_volume", VolumeDiscoverNewDiscsUpdate, "",
				assertComplete(s.Assertions, &vltypes.VolumeDiscoverNewDiscsUpdateResponse{
					Discovered: []string{"test_volume/disc2"},
				}, nil), nil)
		}, time.Hour)
	s.env.ExecuteWorkflow(Volume, state)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &vltypes.VolumeState{
		Discs: []string{"test_volume/disc1", "test_volume/disc2"},
	})
}

func TestVolume(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(volumeTestSuite))
}
