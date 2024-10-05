package vlcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"
)

func TestGetTemporalClient(t *testing.T) {
	t.Parallel()

	// Create a mock Temporal client
	mockClient, err := client.NewLazyClient(client.Options{})
	require.NoError(t, err)

	// Create a context with the Temporal client
	ctx := WithTemporalClient(context.Background(), mockClient)

	// Retrieve the Temporal client from the context
	retrievedClient := GetTemporalClient(ctx)

	// Assert that the retrieved client is the same as the mock client
	assert.Equal(t, mockClient, retrievedClient)
}
