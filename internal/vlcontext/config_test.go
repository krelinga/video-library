package vlcontext

import (
	"context"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	ctx := context.Background()
	config := &vlconfig.Root{
		Volume: &vlconfig.Volume{
			Directory: "/nas/media/Volumes",
		},
	}

	ctxWithConfig := WithConfig(ctx, config)
	retrievedConfig := GetConfig(ctxWithConfig)
	assert.NotNil(t, retrievedConfig)
	assert.Equal(t, config, retrievedConfig)
}
