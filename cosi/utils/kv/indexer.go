package kv

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/numiadata/cosi/utils/pubsub"
)

// Data in the tx_index.db is stored as
// Events:
// key: txhash
// value: tx data

// HeightIndex:
// key: height
// value: txhash

// Loop over all heights, get the tx hash then get the tx data and send to pubsub
// good for tx indexing, but for block events we would still need to got through the state.db (RPC)

const (
	tagKeySeparator = "/"
	txKey           = "tx"
	heightKey       = "height"
)

var i int64

func IndexTxs(ctx context.Context, consumer *pubsub.EventSink, path string, start, end int64) error {

	db, err := newTxIndex(path)
	if err != nil {
		return err
	}
	for i := start; i < end; i++ {
		// get tx hash
		hashes, err := db.getHashes(i)
		if err != nil {
			return err
		}

		results := make([]*abci.TxResult, 0, len(hashes))
		// get tx data
		for _, hash := range hashes {
			events, err := db.getTxEvents(hash) //todo pass in the pubsub streamed
			if err != nil {
				return err
			}
			results = append(results, events)
		}

		// index this blocks txs
		consumer.IndexTxs(results)
	}

	return nil
}

// TxIndexer is an interface for indexing transactions.

// txIndex is the simplest possible indexer, backed by key-value storage (levelDB).
type txIndex struct {
	store dbm.DB
}

// NewTxIndex creates new KV indexer.

func newTxIndex(path string) (*txIndex, error) {
	store, err := dbm.NewGoLevelDBWithOpts("tx_index", path, &opt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	return &txIndex{
		store: store,
	}, nil
}

// getHashes returns the tx hashes for the given height.
func (txi *txIndex) getHashes(height int64) ([][]byte, error) {
	key := startKey(txKey, heightKey, height)
	bz, err := txi.store.Get(key)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, fmt.Errorf("no txs for $%d", height)
	}

	// cheeky way to only print lints every 1000 blocks
	if i == 1000 {
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

func startKey(fields ...interface{}) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write([]byte(fmt.Sprintf("%v", f) + tagKeySeparator))
	}
	return b.Bytes()
}
