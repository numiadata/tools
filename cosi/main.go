package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

func main() {
	fs := flag.NewFlagSet("indexer", flag.ContinueOnError)
	var (
		dir         = fs.String("dir", "xxx", "directory containing tx_index.db LevelDB dir")
		start       = fs.Int64("start", 0, "start height")
		end         = fs.Int64("end", 1000, "end height")
		getHashes   = fs.Bool("get-hashes", true, "get hashes for each height")
		getTxEvents = fs.Bool("get-tx-events", true, "get tx events for each height (implies get-hashes)")
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *getTxEvents {
		*getHashes = true
	}

	log.Printf("dir %s", *dir)
	log.Printf("start %d", *start)
	log.Printf("end %d", *end)
	log.Printf("get-hashes %v", *getHashes)
	log.Printf("get-tx-events %v", *getTxEvents)

	begin := time.Now()
	n, err := noop(context.Background(), *dir, *start, *end, *getHashes, *getTxEvents)
	took := time.Since(begin)
	log.Printf("n=%d", n)
	log.Printf("err=%v", err)
	log.Printf("took=%s", took)
	log.Printf("rate=%.3f heights per second", float64(n)/took.Seconds())
}

func noop(ctx context.Context, dir string, start, end int64, getHashes bool, getTxEvents bool) (count uint64, _ error) {
	txIndex, err := newTxIndex(dir, start, end)
	if err != nil {
		return 0, fmt.Errorf("new tx index: %w", err)
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
		if getHashes {
			hashes, err := txIndex.getHashes(ctx, i)
			if err != nil {
				return count, fmt.Errorf("i=%d: get hashes: %w", i, err)
			}

			if getTxEvents {
				for _, hash := range hashes {
					_, err := txIndex.getTxEvents(hash)
					if err != nil {
						return count, fmt.Errorf("i=%d: get tx events: %w", i, err)
					}
				}
			}
		}
		atomic.AddUint64(&count, 1)
	}

	return count, nil
}

type txIndex struct {
	store dbm.DB
}

func newTxIndex(path string, start, end int64) (*txIndex, error) {
	store, err := dbm.NewGoLevelDBWithOpts("tx_index", path, &opt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	return &txIndex{
		store: store,
	}, nil
}

func (txi *txIndex) getHashes(ctx context.Context, height int64) (map[string][]byte, error) {
	key := heightKey("tx.height", height, height)

	it, err := dbm.IteratePrefix(txi.store, key)
	if err != nil {
		return nil, fmt.Errorf("iterate prefix: %w", err)
	}
	defer it.Close()

	tmpHashes := map[string][]byte{}
	for ; it.Valid(); it.Next() {
		tmpHashes[string(it.Value())] = it.Value()
		if err := ctx.Err(); err != nil {
			break
		}
	}

	return tmpHashes, nil
}

func (txi *txIndex) getTxEvents(hash []byte) (*abci.TxResult, error) {
	rawBytes, err := txi.store.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("store get: %w", err)
	}
	if rawBytes == nil {
		return nil, nil
	}

	var txResult abci.TxResult

	if err := proto.Unmarshal(rawBytes, &txResult); err != nil {
		return nil, fmt.Errorf("unmarshal tx result: %w", err)
	}

	return &txResult, nil
}

func heightKey(fields ...interface{}) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write([]byte(fmt.Sprintf("%v", f) + "/"))
	}
	return b.Bytes()
}
