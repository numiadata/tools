package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"
	abci "github.com/tendermint/tendermint/abci/types"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmstate "github.com/tendermint/tendermint/proto/tendermint/state"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/numiadata/tools/cosi/utils/pubsub"
)

// The state package defines indexing the state.db

func Index(ctx context.Context, consumer *pubsub.EventSink, path string, start, end int64, unsafe bool) error {

	// db, err := newStateStore(path, start, end)
	// if err != nil {
	// 	return err
	// }

	for i := start; i < end; i++ {
		// block, err := db.GetABCIResponses(i)
		// if err != nil {
		// 	// return err
		// }

	}

	// if len(hashes) > 0 {
	// 	results := make([]*abci.TxResult, 0, len(hashes))
	// 	// get tx data
	// 	for _, hash := range hashes {
	// 		events, err := db.getTxEvents(hash)
	// 		if err != nil {
	// 			// return err
	// 		}
	// 		results = append(results, events)
	// 	}

	// 	// index this blocks txs
	// 	consumer.IndexTxs(results, unsafe)
	// }

	return nil
}

type stateStore struct {
	state dbm.DB
	block dbm.DB
}

func newStateStore(path string, start, end int64) (*stateStore, error) {
	state, err := dbm.NewGoLevelDBWithOpts("state", path, &opt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	block, err := dbm.NewGoLevelDBWithOpts("blockstore", path, &opt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	return &stateStore{
		state: state,
		block: block,
	}, nil
}

func (store stateStore) GetEventHeader(height int64) (*types.EventDataNewBlockHeader, *[]abci.TxResult, error) {
	// res, err := store.GetABCIResponses(height)
	// if err != nil {
	// 	return nil, nil, err
	// }

	return nil, nil, nil
}

// GetABCIResponses returns the ABCIResponses for the given height.
func (store stateStore) GetABCIResponses(height int64) (*tmstate.ABCIResponses, error) {

	buf, err := store.state.Get(calcABCIResponsesKey(height))
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("no ABCIResponses for height")
	}

	abciResponses := new(tmstate.ABCIResponses)
	err = abciResponses.Unmarshal(buf)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadABCIResponses: Data has been corrupted or its spec has
                changed: %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return abciResponses, nil
}

func calcABCIResponsesKey(height int64) []byte {
	return []byte(fmt.Sprintf("abciResponsesKey:%v", height))
}

// LoadBlockMeta returns the BlockMeta for the given height.
// If no block is found for the given height, it returns nil.
func (bs *stateStore) LoadBlockMeta(height int64) *types.BlockMeta {
	var pbbm = new(tmproto.BlockMeta)
	bz, err := bs.block.Get(calcBlockMetaKey(height))

	if err != nil {
		panic(err)
	}

	if len(bz) == 0 {
		return nil
	}

	err = proto.Unmarshal(bz, pbbm)
	if err != nil {
		panic(fmt.Errorf("unmarshal to tmproto.BlockMeta: %w", err))
	}

	blockMeta, err := types.BlockMetaFromProto(pbbm)
	if err != nil {
		panic(fmt.Errorf("error from proto blockMeta: %w", err))
	}

	return blockMeta
}

func calcBlockMetaKey(height int64) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}
