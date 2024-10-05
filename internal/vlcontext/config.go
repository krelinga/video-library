package vlcontext

import (
	"context"

	"github.com/krelinga/video-library/internal/vlconfig"
)

type configContextKeyType struct{}
var configContextKey = &configContextKeyType{}

func WithConfig(ctx context.Context, config *vlconfig.Root) context.Context {
	return context.WithValue(ctx, configContextKey, config)
}

func GetConfig(ctx context.Context) *vlconfig.Root {
	return ctx.Value(configContextKey).(*vlconfig.Root)
}