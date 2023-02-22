package state

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	cosdb "github.com/cosmos/cosmos-db"
	"github.com/gogo/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmstate "github.com/tendermint/tendermint/proto/tendermint/state"
	tmstore "github.com/tendermint/tendermint/proto/tendermint/store"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/numiadata/tools/pubsub"
)

// lift count to globabl scope
var count uint64

// The state package defines indexing the state.db
func Index(ctx context.Context, consumer *pubsub.EventSink, path, db string, start, end int64, unsafe bool) error {

	statedb, err := newStateStore(path, db)
	if err != nil {
		return fmt.Errorf("new stateStore: %w", err)
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin := time.Now()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				delta := time.Since(begin)
				count := atomic.LoadUint64(&count)
				rate := float64(count) / float64(delta.Seconds())
				log.Printf("+%ds count=%d rate=%.3f/s", int(delta.Seconds()), count, rate)
			}
		}
	}()

	for i := start; i < end; i++ {

		fmt.Println("height", i)
		// indexing blocks
		res, err := statedb.getABCIResponses(i)
		if err != nil {
			return fmt.Errorf("i=%d: get abciresponses: %w", i, err)
		}

		b := statedb.loadBlock(i)

		e := tmtypes.EventDataNewBlockHeader{
			ResultBeginBlock: *res.BeginBlock,
			ResultEndBlock:   *res.EndBlock,
			Header:           b.Header,
			NumTxs:           int64(len(b.Data.Txs)),
		}

		consumer.IndexBlock(e, true)

		// indexing txs

		var batch = make([]*pubsub.TxResult, 0, len(b.Data.Txs))
		if e.NumTxs > 0 {

			for i := range b.Data.Txs {
				tr := &pubsub.TxResult{
					Height:         b.Height,
					Index:          uint32(i),
					Tx:             b.Data.Txs[i],
					Result:         *(res.DeliverTxs[i]),
					BlockTimestamp: b.Time,
				}
				batch = append(batch, tr)

			}

			consumer.IndexTxs(batch, true)
		}

		atomic.AddUint64(&count, 1)
	}

	return nil
}

func GetBaseHeight(path, db string) (int64, int64, error) {
	statedb, err := newStateStore(path, db)
	if err != nil {
		return 0, 0, fmt.Errorf("new stateStore: %w", err)
	}

	base, height := statedb.loadBlockStoreState()

	return base, height, nil
}

func ForceCompact(path string) error {

	state, err := dbm.NewGoLevelDBWithOpts("state", path, &opt.Options{Filter: filter.NewBloomFilter(10)})
	if err != nil {
		return err
	}

	block, err := dbm.NewGoLevelDBWithOpts("blockstore", path, &opt.Options{Filter: filter.NewBloomFilter(10)})
	if err != nil {
		return err
	}

	fmt.Println("compacting state db")
	if err := state.ForceCompact(nil, nil); err != nil {
		return fmt.Errorf("force compact state: %w", err)
	}

	fmt.Println("compacting block db")
	err = block.ForceCompact(nil, nil)

	return fmt.Errorf("force compact blockstore: %w", err)
}

type stateStore struct {
	state cosdb.DB
	block cosdb.DB
}

func newStateStore(path string, db string) (*stateStore, error) {
	var (
		state cosdb.DB
		block cosdb.DB
	)

	fmt.Println(db)
	if db == "goleveldb" {
		st, err := cosdb.NewGoLevelDBWithOpts("state", path, &opt.Options{ReadOnly: true})
		if err != nil {
			return nil, err
		}

		bl, err := cosdb.NewGoLevelDBWithOpts("blockstore", path, &opt.Options{ReadOnly: true})
		if err != nil {
			return nil, err
		}

		state = st
		block = bl
	} else if db == "pebbledb" {
		fmt.Println(1)
		st, err := cosdb.NewDB("state", cosdb.PebbleDBBackend, path)
		if err != nil {
			return nil, err
		}

		bl, err := cosdb.NewDB("blockstore", cosdb.PebbleDBBackend, path)
		if err != nil {
			return nil, err
		}

		state = st
		block = bl
	}

	return &stateStore{
		state: state,
		block: block,
	}, nil
}

// GetABCIResponses returns the ABCIResponses for the given height.
func (store stateStore) getABCIResponses(height int64) (*tmstate.ABCIResponses, error) {
	fmt.Println(height)
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
	fmt.Println(height)
	return []byte(fmt.Sprintf("abciResponsesKey:%v", height))
}

// LoadBlockMeta returns the BlockMeta for the given height.
// If no block is found for the given height, it returns nil.
func (bs *stateStore) loadBlockMeta(height int64) (*types.BlockMeta, error) {
	var pbbm = new(tmproto.BlockMeta)
	bz, err := bs.block.Get(calcBlockMetaKey(height))

	if err != nil {
		panic(err)
	}

	if len(bz) == 0 {
		return nil, nil
	}

	err = proto.Unmarshal(bz, pbbm)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to tmproto.BlockMeta: %w", err)
	}

	blockMeta, err := types.BlockMetaFromProto(pbbm)
	if err != nil {
		return nil, fmt.Errorf("error from proto blockMeta: %w", err)
	}

	return blockMeta, nil
}

func calcBlockMetaKey(height int64) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}

func (bs *stateStore) loadBlock(height int64) *types.Block {
	blockMeta, err := bs.loadBlockMeta(height)
	if err != nil {
		panic(err)
	}
	if blockMeta == nil {
		return nil
	}

	pbb := new(tmproto.Block)
	buf := []byte{}
	for i := 0; i < int(blockMeta.BlockID.PartSetHeader.Total); i++ {
		part := bs.loadBlockPart(height, i)
		// If the part is missing (e.g. since it has been deleted after we
		// loaded the block meta) we consider the whole block to be missing.
		if part == nil {
			return nil
		}
		buf = append(buf, part.Bytes...)
	}
	err = proto.Unmarshal(buf, pbb)
	if err != nil {
		// NOTE: The existence of meta should imply the existence of the
		// block. So, make sure meta is only saved after blocks are saved.
		panic(fmt.Sprintf("Error reading block: %v", err))
	}

	block, err := types.BlockFromProto(pbb)
	if err != nil {
		panic(fmt.Errorf("error from proto block: %w", err))
	}

	return block
}

func (bs *stateStore) loadBlockPart(height int64, index int) *types.Part {
	var pbpart = new(tmproto.Part)

	bz, err := bs.block.Get(calcBlockPartKey(height, index))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil
	}

	err = proto.Unmarshal(bz, pbpart)
	if err != nil {
		panic(fmt.Errorf("unmarshal to tmproto.Part failed: %w", err))
	}
	part, err := types.PartFromProto(pbpart)
	if err != nil {
		panic(fmt.Sprintf("Error reading block part: %v", err))
	}

	return part
}

func calcBlockPartKey(height int64, partIndex int) []byte {
	return []byte(fmt.Sprintf("P:%v:%v", height, partIndex))
}

// LoadBlockStoreState returns the BlockStoreState as loaded from disk.
// If no BlockStoreState was previously persisted, it returns the zero value.
var blockStoreKey = []byte("blockStore")

func (bs *stateStore) loadBlockStoreState() (base int64, height int64) {
	bytes, err := bs.block.Get(blockStoreKey)
	if err != nil {
		panic(err)
	}

	if len(bytes) == 0 {
		return 0, 0
	}

	var bsj tmstore.BlockStoreState
	if err := proto.Unmarshal(bytes, &bsj); err != nil {
		panic(fmt.Sprintf("Could not unmarshal bytes: %X", bytes))
	}

	// Backwards compatibility with persisted data from before Base existed.
	if bsj.Height > 0 && bsj.Base == 0 {
		bsj.Base = 1
	}

	return bsj.Base, bsj.Height
}
