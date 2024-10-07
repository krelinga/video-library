package vltemp

import (
	"context"
	"testing"

	"github.com/krelinga/video-library/internal/vlconfig"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/stretchr/testify/assert"
)

func TestDiscPath(t *testing.T) {
	ctx := vlcontext.WithConfig(context.Background(), &vlconfig.Root{Volume: &vlconfig.Volume{Directory: "/mocked/path"}})
	gotPath, err := DiscPath(ctx, DiscWfId("volume1/disc1"))
	assert.NoError(t, err)
	assert.Equal(t, "/mocked/path/volume1/disc1", gotPath)
}
