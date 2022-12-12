package io

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsFileComplete determines if a state streamed file is completely written. The
// Cosmos SDK state streamer that streams to files prefixes the files with the
// length encoded size of the data. The length prefix is scanned to determine if
// all expected bytes are written.
func IsFileComplete(filePath string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	defer f.Close()

	buf := make([]byte, 8)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(buf, 8)

	ok := scanner.Scan()
	if ok {
		return false, errors.New("expected scanner to stop")
	}

	length := sdk.BigEndianToUint64(buf)
	if length == 0 {
		return false, errors.New("expected positive length prefix size")
	}

	fInfo, err := f.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
	}

	return length == uint64(fInfo.Size())-8, nil
}
