package vltemp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
)

type volumeTestSuite struct {
	testSuite
}

func (s *volumeTestSuite) TestFreshlyCreated() {
	s.env.OnActivity(actVolumeMkDir, mock.Anything, VolumeWfId("test_volume")).Return(nil)
	s.env.SetStartWorkflowOptions(client.StartWorkflowOptions{
		ID: "test_volume",
	})
	s.env.ExecuteWorkflow(VolumeWF, nil)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assertContinuedWithState(s.Assertions, err, &VolumeWFState{})
}

func (s *volumeTestSuite) TestDiscoverNewDiscs() {
	s.env.OnActivity(actVolumeReadDiscNames, mock.Anything, VolumeWfId("test_volume")).Return([]string{"disc1", "disc2"}, nil)
	s.env.OnActivity(actVolumeBootstrapDisc, mock.Anything, VolumeWfId("test_volume"), "disc2").Return(DiscWfId("test_volume/disc2"), nil)
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
func TestVolumeWfId(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Valid ID", "valid-volume-id", false},
		{"Invalid ID with slashes", "invalid/volume/id", true},
		{"Invalid ID with colons", "invalid:volume:id", true},
		{"Empty ID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewVolumeWfId(tt.input)
			if tt.expectErr {
				assert.ErrorIs(t, err, ErrInvalidWorkflowId)
				assert.Panics(t, func() { id.Name() })
			} else {
				assert.NoError(t, err)
				assert.Equal(t, VolumeWfId(tt.input), id)
			}
		})
	}
}
