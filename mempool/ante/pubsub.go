package ante

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	credsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

	// AttrKeyChainID     = "chain_id"
	// AttrKeyBlockHeight = "block_height"
	// AttrKeyTxHash      = "tx_hash"
	// AttrKeyTxCount     = "tx_count"

	// MsgType            = "message_type"
	// MsgTypeBeginBlock  = "begin_block"
	// MsgTypeBlockHeader = "header"
	// MsgTypeEndBlock    = "end_block"
	// MsgTypeTxResult    = "tx_result"
	// MsgTypeTxEvents    = "tx_events"
	// MsgTypeTxCount     = "tx_count"
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
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPubSubDecorator(projectID, topic string) PubSubDecorator {
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
		client: psClient,
		topic:  psTopic,
	}
}

func (d PubSubDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	panic("not implemented!")
}
