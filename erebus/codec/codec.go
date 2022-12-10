package codec

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

var (
	iRegistry = codectypes.NewInterfaceRegistry()
	cdc       = codec.NewProtoCodec(iRegistry)
)

// ParseDataFile attempts to read in a streamed data file, which contains one or
// more length-prefixed StoreKVPair records and returns a slice of StoreKVPair
// records decoded. An error is returned upon failure.
func ParseDataFile(f string) ([]storetypes.StoreKVPair, error) {
	bz, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	segments, err := segmentBytes(bz[8:]) // strip the useless file length prefix prefix
	if err != nil {
		return nil, fmt.Errorf("failed to parse data file: %w", err)
	}

	pairs := make([]storetypes.StoreKVPair, 0, len(segments))
	for _, s := range segments {
		var pair storetypes.StoreKVPair
		if err := cdc.Unmarshal(s, &pair); err != nil {
			return nil, fmt.Errorf("failed to decode StoreKVPair from segment: %w", err)
		}

		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// ParseMetaFile attempts to read in a streamed metadata file, which contains
// a single Protobuf length-prefixed encoded BlockMetadata and return a decoded
// BlockMetadata object. An error is returned upon failure.
func ParseMetaFile(f string) (storetypes.BlockMetadata, error) {
	bz, err := os.ReadFile(f)
	if err != nil {
		return storetypes.BlockMetadata{}, fmt.Errorf("failed to read meta file: %w", err)
	}

	// strip the useless file length prefix prefix
	var metadata storetypes.BlockMetadata
	if err := cdc.Unmarshal(bz[8:], &metadata); err != nil {
		return storetypes.BlockMetadata{}, fmt.Errorf("failed to decode BlockMetadata: %w", err)
	}

	return metadata, nil
}

// segmentBytes returns all of the Protobuf encoded messages contained in the
// slice of byte slices. Each byte slice corresponds to a Protobuf encoded
// message. Each message has it's length prefix removed.
func segmentBytes(bz []byte) ([][]byte, error) {
	var err error

	segments := make([][]byte, 0)
	for len(bz) > 0 {
		var segment []byte
		segment, bz, err = getHeadSegment(bz)
		if err != nil {
			return nil, err
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// getHeadSegment returns the bytes for the leading protobuf object in the byte
// slice (removing the length prefix) and returns the remainder of the byte
// slice.
func getHeadSegment(bz []byte) ([]byte, []byte, error) {
	size, prefixSize := binary.Uvarint(bz)
	if prefixSize < 0 {
		return nil, nil, fmt.Errorf("invalid number of bytes read from length-prefixed encoding: %d", prefixSize)
	}

	if size > uint64(len(bz)-prefixSize) {
		return nil, nil, fmt.Errorf("not enough bytes to read; want: %v, got: %v", size, len(bz)-prefixSize)
	}

	return bz[prefixSize:(uint64(prefixSize) + size)], bz[uint64(prefixSize)+size:], nil
}
