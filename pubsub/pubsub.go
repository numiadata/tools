package pubsub

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/tendermint/tendermint/types"
)

const (
	credsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

	//todo: migrate this to flags
	projectIDEnvVar = "PROJECT_ID"
	topicEnvVar     = "TOPIC"
	chainIDEnvVar   = "CHAIN_ID"

	AttrKeyChainID     = "chain_id"
	AttrKeyBlockHeight = "block_height"
	AttrKeyTxHash      = "tx_hash"

	MsgType            = "message_type"
	MsgTypeBeginBlock  = "begin_block"
	MsgTypeBlockHeader = "header"
	MsgTypeEndBlock    = "end_block"
	MsgTypeTxResult    = "tx_result"
	MsgTypeTxEvents    = "tx_events"
)

var jsonpbMarshaller = jsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: true,
}

type EventSink struct {
	client  *pubsub.Client
	topic   *pubsub.Topic
	chainID string
}

func NewEventSink() (*EventSink, error) {

	if s := os.Getenv(credsEnvVar); len(s) == 0 {
		return nil, fmt.Errorf("missing '%s' environment variable", credsEnvVar)
	}

	projectID := os.Getenv(projectIDEnvVar)
	if len(projectID) == 0 {
		return nil, fmt.Errorf("missing '%s' environment variable", projectIDEnvVar)
	}

	topic := os.Getenv(topicEnvVar)
	if len(topic) == 0 {
		return nil, fmt.Errorf("missing '%s' environment variable", topicEnvVar)
	}

	chainID := os.Getenv(chainIDEnvVar)
	if len(chainID) == 0 {
		return nil, fmt.Errorf("missing '%s' environment variable", chainIDEnvVar)
	}

	fmt.Println("started pubsub")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a Google Cloud Pubsub client: %w", err)
	}

	// attempt to get the topic. If that fails, we attempt to create it
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check for topic '%s': %w", topic, err)
	}

	if !ok {
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		t, err = c.CreateTopic(ctx, topic)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic '%s': %w", topic, err)
		}
	}

	return &EventSink{
		client:  c,
		topic:   t,
		chainID: chainID,
	}, nil
}

// IndexBlock attempts to index ResultBeginBlock and ResultEndBlock block events,
// where the data in each pubsub message includes the JSON encoding of the
// respective type. Additional attributes are included such as the block height,
// the chain ID and the message type. An error is returned if any encoding error
// is detected or if any message fails to publish.

func (es *EventSink) IndexBlock(h types.EventDataNewBlockHeader, unsafe bool) error {
	buf := new(bytes.Buffer)
	blockHeightStr := strconv.Itoa(int(h.Header.Height))

	var results []*pubsub.PublishResult

	// publish block header
	if err := jsonpbMarshaller.Marshal(buf, h.Header.ToProto()); err != nil {
		return fmt.Errorf("failed to JSON marshal Header: %w", err)
	}

	res := es.topic.Publish(
		context.Background(), // NOTE: contexts aren't used in Publish
		&pubsub.Message{
			Data: buf.Bytes(),
			Attributes: map[string]string{
				MsgType:            MsgTypeBlockHeader,
				AttrKeyChainID:     es.chainID,
				AttrKeyBlockHeight: blockHeightStr,
			},
		},
	)
	results = append(results, res)

	buf1 := new(bytes.Buffer)

	// publish BeginBlock Events
	if err := jsonpbMarshaller.Marshal(buf1, &h.ResultBeginBlock); err != nil {
		return fmt.Errorf("failed to JSON marshal ResultBeginBlock: %w", err)
	}

	res = es.topic.Publish(
		context.Background(), // NOTE: contexts aren't used in Publish
		&pubsub.Message{
			Data: buf1.Bytes(),
			Attributes: map[string]string{
				MsgType:            MsgTypeBeginBlock,
				AttrKeyChainID:     es.chainID,
				AttrKeyBlockHeight: blockHeightStr,
			},
		},
	)
	results = append(results, res)

	// buf.Reset() // reset buffer prior to next Marshal call
	buf2 := new(bytes.Buffer)

	// publish EndBlock Events
	if err := jsonpbMarshaller.Marshal(buf2, &h.ResultEndBlock); err != nil {
		return fmt.Errorf("failed to JSON marshal ResultEndBlock: %w", err)
	}

	res = es.topic.Publish(
		context.Background(), // NOTE: contexts aren't used in Publish
		&pubsub.Message{
			Data: buf2.Bytes(),
			Attributes: map[string]string{
				MsgType:            MsgTypeEndBlock,
				AttrKeyChainID:     es.chainID,
				AttrKeyBlockHeight: blockHeightStr,
			},
		},
	)
	results = append(results, res)

	if !unsafe {
		fmt.Println("safe")
		// wait for all messages to be be sent (or failed to be sent) to the server
		for _, r := range results {
			if _, err := r.Get(context.Background()); err != nil {
				return fmt.Errorf("failed to publish pubsub message: %w", err)
			}
		}
	}
	return nil
}

func (es *EventSink) IndexTxs(txrs []*TxResult, unsafe bool) error {
	results := make([]*pubsub.PublishResult, len(txrs))
	for i, txr := range txrs {
		buf := new(bytes.Buffer)
		// buf.Reset() // reset buffer prior to next Marshal call

		blockHeightStr := strconv.Itoa(int(txr.Height))
		txHash := fmt.Sprintf("%X", types.Tx(txr.Tx).Hash())

		if err := jsonpbMarshaller.Marshal(buf, txr); err != nil {
			return fmt.Errorf("failed to JSON marshal TxResult: %w", err)
		}

		res := es.topic.Publish(
			context.Background(), // NOTE: contexts aren't used in Publish
			&pubsub.Message{
				Data: buf.Bytes(),
				Attributes: map[string]string{
					MsgType:            MsgTypeTxResult,
					AttrKeyChainID:     es.chainID,
					AttrKeyBlockHeight: blockHeightStr,
					AttrKeyTxHash:      txHash,
				},
			},
		)

		results[i] = res
	}

	if !unsafe {
		fmt.Println("safe")
		// wait for all messages to be be sent (or failed to be sent) to the server
		for _, r := range results {
			if _, err := r.Get(context.Background()); err != nil {
				return fmt.Errorf("failed to publish pubsub message: %w", err)
			}
		}
	}

	return nil
}
