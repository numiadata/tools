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
