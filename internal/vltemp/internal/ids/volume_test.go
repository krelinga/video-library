package ids

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVolumeWfId(t *testing.T) {
	tests := []struct {
		name       string
		volumeName string
		expectErr  bool
	}{
		{
			name:       "valid volume name",
			volumeName: "test_volume",
			expectErr:  false,
		},
		{
			name:       "empty volume name",
			volumeName: "",
			expectErr:  true,
		},
		{
			name:       "invalid volume name with slash",
			volumeName: "invalid/volume",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewVolumeWfId(tt.volumeName)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, id)
				assert.Equal(t, ErrInvalidWorkflowId, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, id)
				assert.Equal(t, tt.volumeName, id.Name())
			}
		})
	}
}
