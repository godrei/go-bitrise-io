package appledevconn_test

import (
	"testing"

	"github.com/godrei/go-bitrise-io/appledevconn"
	"github.com/stretchr/testify/require"
)

func TestEnsureConnection(t *testing.T) {
	require.NoError(t, appledevconn.EnsureConnection())
}
