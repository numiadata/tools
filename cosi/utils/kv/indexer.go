package kv

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/google/orderedcode"
	"github.com/syndtr/goleveldb/leveldb/opt"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/alitto/pond"
	"github.com/numiadata/tools/cosi/utils/pubsub"
)

// Data in the tx_index.db is stored as
// Events:
// key: txhash
// value: tx data

// HeightIndex:
// key: height
// value: txhash

const (
	tagKeySeparator = "/"
	txKey           = "tx.height"

	blockPrefix = "block_events"
)

func IndexTxs(ctx context.Context, consumer *pubsub.EventSink, path string, start, end int64, unsafe bool) error {

	db, err := newTxIndex(path, start, end)
	if err != nil {
		return err
	}

	pool := pond.New(250, 2500)

	//TODO see if this can be concurrent
	for i := start; i < end; i++ {
		n := i
		pool.Submit(func() {
			Index(db, ctx, consumer, n, unsafe)
		})

	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()

	return nil
}

func Index(db *txIndex, ctx context.Context, consumer *pubsub.EventSink, i int64, unsafe bool) {
	// get tx hash
	hashes, err := db.getHashes(ctx, i)
	if err != nil {
		// return err
	}

	if len(hashes) == 0 {
		fmt.Println("no txs for height", i)

	}
	fmt.Println("height ", i)

	if len(hashes) > 0 {
		results := make([]*abci.TxResult, 0, len(hashes))
		// get tx data
		for _, hash := range hashes {
			events, err := db.getTxEvents(hash)
			if err != nil {
				// return err
			}
			results = append(results, events)
		}

		// index this blocks txs
		consumer.IndexTxs(results, unsafe)
	}

	// return nil
}

// TxIndexer is an implementation for indexing transactions.

// txIndex is the simplest possible indexer, backed by key-value storage (levelDB).
type txIndex struct {
	store dbm.DB
}

// NewTxIndex creates new KV indexer.

func newTxIndex(path string, start, end int64) (*txIndex, error) {
	store, err := dbm.NewGoLevelDBWithOpts("tx_index", path, &opt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	// sKey := heightKey(txKey, heightKey, start)
	// endKey := heightKey(txKey, heightKey, end)
	// store.Iterator(sKey, endKey)

	return &txIndex{
		store: store,
	}, nil
}

// getHashes returns the tx hashes for the given height.
func (txi *txIndex) getHashes(ctx context.Context, height int64) (map[string][]byte, error) {
	key := heightKey(txKey, height, height)

	it, err := dbm.IteratePrefix(txi.store, key)
	if err != nil {
		return nil, err
	}

	defer it.Close()

	tmpHashes := make(map[string][]byte)

	for ; it.Valid(); it.Next() {
		tmpHashes[string(it.Value())] = it.Value()

		// Potentially exit early.
		select {
		case <-ctx.Done():
			break
		default:
		}
	}

	// bz, err := txi.store.Get(key)
	// if err != nil {
	// 	return nil, err
	// }

	if len(tmpHashes) == 0 {
		return nil, nil
	}

	return tmpHashes, nil
}

// getTxEvents returns the events for the given tx hash and sends it via pubsub to a listener
func (txi *txIndex) getTxEvents(hash []byte) (*abci.TxResult, error) {

	rawBytes, err := txi.store.Get(hash)
	if err != nil {
		panic(err)
	}
	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(abci.TxResult)
	err = proto.Unmarshal(rawBytes, txResult)

	return txResult, err
}

func heightKey(fields ...interface{}) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write([]byte(fmt.Sprintf("%v", f) + tagKeySeparator))
	}
	return b.Bytes()
}

// BlockIndexer is an implementation for indexing block events.

type blockerIndexer struct {
	store dbm.DB
}

func New(store dbm.DB) *blockerIndexer {
	prefixStore := dbm.NewPrefixDB(store, []byte(blockPrefix))
	return &blockerIndexer{
		store: prefixStore,
	}
}

// Has returns true if the given height has been indexed. An error is returned
// upon database query failure.
func (idx *blockerIndexer) Has(height int64) (bool, error) {
	key, err := blockHeightKey(height)
	if err != nil {
		return false, fmt.Errorf("failed to create block height index key: %w", err)
	}

	return idx.store.Has(key)
}

func blockHeightKey(height int64) ([]byte, error) {
	return orderedcode.Append(
		nil,
		types.BlockHeightKey,
		height,
	)
}
