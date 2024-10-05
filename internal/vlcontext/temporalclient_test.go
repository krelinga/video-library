package vlcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
)

func TestGetTemporalClient(t *testing.T) {
	t.Parallel()

	// Create a mock Temporal client
	mockClient := new(client.Client)

	// Create a context with the Temporal client
	ctx := WithTemporalClient(context.Background(), mockClient)

	// Retrieve the Temporal client from the context
	retrievedClient := GetTemporalClient(ctx)

	// Assert that the retrieved client is the same as the mock client
	assert.Equal(t, mockClient, retrievedClient)
}
