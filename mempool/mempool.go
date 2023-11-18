package mempool

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"
)

var _ sdkmempool.Mempool = (*PubSubMempool)(nil)

// PubSubMempool defines an SDK mempool, which internally extends a no-op mempool
// by overriding the Insert method which solely emits Google Cloud pubsub events
// about transactions entering the mempool.
type PubSubMempool struct {
	sdkmempool.NoOpMempool

	logger    log.Logger
	chainID   string
	nodeID    string
	txEncoder sdk.TxEncoder
	client    *pubsub.Client
	topic     *pubsub.Topic
	sync      bool // sync defines if we should wait for all pubsub results to complete prior to returning
}

func NewPubSubMempool(logger log.Logger, txEncoder sdk.TxEncoder, chainID, nodeID, projectID, topic string, sync bool) *PubSubMempool {
	if s := os.Getenv(CredsEnvVar); len(s) == 0 {
		panic(fmt.Errorf("missing '%s' environment variable", CredsEnvVar))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	psClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(fmt.Errorf("failed to create a Google Cloud Pubsub client: %w", err))
	}

	// Attempt to get the topic. If that fails, we attempt to create it.
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	psTopic := psClient.Topic(topic)

	topicExists, err := psTopic.Exists(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to check for topic '%s': %w", topic, err))
	}

	if !topicExists {
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		psTopic, err = psClient.CreateTopic(ctx, topic)
		if err != nil {
			panic(fmt.Errorf("failed to create topic '%s': %w", topic, err))
		}
	}

	return &PubSubMempool{
		NoOpMempool: sdkmempool.NoOpMempool{},
		logger:      logger.With("module", "pubsub_mempool"),
		chainID:     chainID,
		nodeID:      nodeID,
		txEncoder:   txEncoder,
		client:      psClient,
		topic:       psTopic,
		sync:        sync,
	}
}

func (mp *PubSubMempool) Insert(ctx context.Context, tx sdk.Tx) error {
	txBz, err := mp.txEncoder(tx)
	if err != nil {
		mp.logger.Error("failed to encode tx", "err", err)
		return nil
	}

	txHashStr := fmt.Sprintf("%X", sha256.Sum256(txBz))
	msgs := tx.GetMsgs()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	results := make([]*pubsub.PublishResult, len(msgs))
	for i, msg := range msgs {
		results[i] = mp.topic.Publish(
			context.Background(),
			&pubsub.Message{
				Data: nil, // TODO(bez): Should we publish the entire tx or just the message?
				Attributes: map[string]string{
					AttrKeyMsgType:   MsgTypeCheckTxMsg,
					AttrKeyChainID:   mp.chainID,
					AttrKeyTxHash:    txHashStr,
					AttrKeyTimestamp: timestamp,
					AttrKeyNodeID:    mp.nodeID,
					AttrKeyMsgSigner: msg.GetSigners()[0].String(),
					AttrKeyTxMsgType: sdk.MsgTypeURL(msg),
				},
			},
		)
	}

	if mp.sync {
		// wait for all messages to be be sent (or failed to be sent) to the server
		for _, r := range results {
			if _, err := r.Get(context.Background()); err != nil {
				mp.logger.Error("failed to publish pubsub message", "err", err)
			}
		}
	}

	return nil
}
