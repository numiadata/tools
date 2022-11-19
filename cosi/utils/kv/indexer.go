package kv

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

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
)

var i int64

func IndexTxs(ctx context.Context, consumer *pubsub.EventSink, path string, start, end int64) error {

	db, err := newTxIndex(path, start, end)
	if err != nil {
		return err
	}
	//TODO see if this can be concurrent
	for i := start; i < end; i++ {
		// get tx hash
		hashes, err := db.getHashes(ctx, i)
		if err != nil {
			return err
		}

		if len(hashes) == 0 {
			fmt.Println("no txs for height", i)
			continue
		}

		if len(hashes) > 0 {
			results := make([]*abci.TxResult, 0, len(hashes))
			// get tx data
			for _, hash := range hashes {
				events, err := db.getTxEvents(hash)
				if err != nil {
					return err
				}
				results = append(results, events)
			}

			// index this blocks txs
			consumer.IndexTxs(results)
		}
	}

	return nil
}

// TxIndexer is an interface for indexing transactions.

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
func (txi *txIndex) getHashes(ctx context.Context, height int64) ([][]byte, error) {
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

	if len(tmpHashes) > 0 {
		fmt.Println(tmpHashes, "Hashes\n\n\n")
	}

	bz, err := txi.store.Get(key)
	if err != nil {
		return nil, err
	}

	if bz == nil {
		return nil, nil
	}

	// cheeky way to only print lints every 1000 blocks
	if i == 500 {
		fmt.Println(height)
		i = 0
	} else {
		i++
	}

	return nil, nil
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
