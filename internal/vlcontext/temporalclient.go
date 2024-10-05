package vlcontext

import (
	"context"

	"go.temporal.io/sdk/client"
)

type temporalClientKeyType struct{}

var temporalClientKey = &temporalClientKeyType{}

func WithTemporalClient(ctx context.Context, client *client.Client) context.Context {
	return context.WithValue(ctx, temporalClientKey, client)
}

func GetTemporalClient(ctx context.Context) *client.Client {
	return ctx.Value(temporalClientKey).(*client.Client)
}
