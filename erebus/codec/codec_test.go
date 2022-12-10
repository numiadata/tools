package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDataFile(t *testing.T) {
	records, err := ParseDataFile("../examples/state_streaming/block-271-data")
	require.NoError(t, err)
	require.NotEmpty(t, records)
}

func TestParseMetaFile(t *testing.T) {
	metadata, err := ParseMetaFile("../examples/state_streaming/block-288-meta")
	require.NoError(t, err)
	require.NotNil(t, metadata.RequestBeginBlock)
	require.NotNil(t, metadata.ResponseBeginBlock)
	require.NotNil(t, metadata.RequestEndBlock)
	require.NotNil(t, metadata.ResponseEndBlock)
	require.NotEmpty(t, metadata.DeliverTxs)
}
