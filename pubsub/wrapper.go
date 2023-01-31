package pubsub

import (
	"errors"

	"github.com/tendermint/tendermint/types"
)

// BlockIndexer implements a wrapper around the Pubsub sink and supports block
// indexing by implementing the indexer.BlockIndexer interface.
type BlockIndexer struct {
	sink *EventSink
}

func NewBlockIndexer(sink *EventSink) *BlockIndexer {
	return &BlockIndexer{sink: sink}
}

func (bi *BlockIndexer) Has(_ int64) (bool, error) {
	return false, errors.New("the Has method is not supported for the Pubsub indexer")
}

func (bi *BlockIndexer) Index(block types.EventDataNewBlockHeader) error {
	return bi.sink.IndexBlock(block, false)
}

// TxIndexer implements a wrapper around the Pubsub sink and supports tx
// indexing by implementing the txindex.TxIndexer interface.
type TxIndexer struct {
	sink *EventSink
}

func NewTxIndexer(sink *EventSink) *TxIndexer {
	return &TxIndexer{sink: sink}
}

func (ti *TxIndexer) AddBatch(batch *Batch) error {
	ops := make([]*TxResult, len(batch.Ops))
	for i, tx := range batch.Ops {
		ops[i] = &TxResult{Tx: tx.Tx, Height: tx.Height, Index: tx.Index, Result: tx.Result} // TODO: missing timestamp
	}

	return ti.sink.IndexTxs(ops, false)
}

type Batch struct {
	Ops []*TxResult
}

func (ti *TxIndexer) Index(txr *TxResult) error {
	op := &TxResult{Tx: txr.Tx, Height: txr.Height, Index: txr.Index, Result: txr.Result} // TODO: missing timestamp

	return ti.sink.IndexTxs([]*TxResult{op}, false)
}

func (ti *TxIndexer) Get(hash []byte) (*TxResult, error) {
	return nil, errors.New("the Get method is not supported for the Pubsub indexer")
}
