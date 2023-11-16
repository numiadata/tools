package ante

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	credsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

	AttrKeyMsgType   = "message_type"
	AttrKeyChainID   = "chain_id"
	AttrKeyTxHash    = "tx_hash"
	AttrKeyTimestamp = "timestamp"
	AttrKeyNodeID    = "node_id"
	AttrKeyMsgSigner = "msg_signer"
	AttrKeyTxMsgType = "tx_msg_type"

	MsgTypeCheckTxMsg = "check_tx_msg"
)

// PubSubDecorator defines an AnteHandler decorator that is responsible for emitting
// transaction context events to a Google Cloud PubSub topic. The decorator publishes
// events during CheckTx only and allows clients to trace the transaction lifecycle
// from CheckTx, which includes inclusion into a mempool, and the transition being
// included in a block.
//
// The decorator should be added to the end of the AnteHandler chain to ensure
// we only publish upon a successful CheckTx execution.
//
// Note, operators must ensure the <GOOGLE_APPLICATION_CREDENTIALS> environment
// variable is set to the location of their creds file.
type PubSubDecorator struct {
	logger log.Logger
	nodeID string
	client *pubsub.Client
	topic  *pubsub.Topic
	sync   bool // sync defines if we should wait for all pubsub results to complete prior to	returning
}

func NewPubSubDecorator(logger log.Logger, nodeID, projectID, topic string, sync bool) PubSubDecorator {
	if s := os.Getenv(credsEnvVar); len(s) == 0 {
		panic(fmt.Errorf("missing '%s' environment variable", credsEnvVar))
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

	return PubSubDecorator{
		logger: logger.With("module", "ante_pubsub"),
		nodeID: nodeID,
		client: psClient,
		topic:  psTopic,
		sync:   sync,
	}
}

func (d PubSubDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		txHashStr := fmt.Sprintf("%X", sha256.Sum256(ctx.TxBytes()))
		msgs := tx.GetMsgs()

		results := make([]*pubsub.PublishResult, len(msgs))
		for i, msg := range msgs {
			results[i] = d.topic.Publish(
				context.Background(),
				&pubsub.Message{
					Data: nil, // TODO(bez): Should we publish the entire tx or just the message?
					Attributes: map[string]string{
						AttrKeyMsgType:   MsgTypeCheckTxMsg,
						AttrKeyChainID:   ctx.ChainID(),
						AttrKeyTxHash:    txHashStr,
						AttrKeyTimestamp: time.Now().UTC().Format(time.RFC3339),
						AttrKeyNodeID:    d.nodeID,
						AttrKeyMsgSigner: msg.GetSigners()[0].String(),
						AttrKeyTxMsgType: sdk.MsgTypeURL(msg),
					},
				},
			)
		}

		if d.sync {
			// wait for all messages to be be sent (or failed to be sent) to the server
			for _, r := range results {
				if _, err := r.Get(context.Background()); err != nil {
					d.logger.Error("failed to publish pubsub message", "err", err)
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
