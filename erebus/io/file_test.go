package io_test

import (
	"testing"

	"github.com/numiadata/tools/erebus/io"
	"github.com/stretchr/testify/require"
)

func TestIsFileComplete(t *testing.T) {
	ok, err := io.IsFileComplete("../examples/state_streaming/block-271-data")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = io.IsFileComplete("../examples/state_streaming/block-271-data-incomplete")
	require.NoError(t, err)
	require.False(t, ok)
}
