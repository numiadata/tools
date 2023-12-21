# mempool

The `mempool` package provides components to emit Google Cloud PubSub messages
upon ingress of Cosmos SDK messages into the mempool. This provides the ability
for clients to infer context about the lifecycle of a transaction, such as when
it enters the mempool and when it's finally included in a block (if at all).

## Usage

The are currently two methods of integration.

### AnteHandler

We provide a Cosmos SDK AnteHandler decorator that a chain can simply inject into
their existing AnteHandler chain. See `PubSubDecorator` for more details and required
arguments for successful integration.

To integrate into an application, simply add the decorator to the existing chain:

```go
import (
  // ...
  antepubsub "github.com/numiadata/tools/mempool/ante
)

func newAnteDecoratorChain(logger log.Logger, opts HandlerOptions) []sdk.AnteDecorator {
  return []sdk.AnteDecorator{
    // ...
    antepubsub.NewPubSubDecorator(logger log.Logger, "<nodeID>", "<projectID>", "<topic>", false),
  }
}
```

> Note, it's best to add the decorator to the very end of the chain as to not emit
> pubsub events for transactions that could fail CheckTx.

### Mempool

We provide an SDK mempool, which internally extends a provided SDK mempool and
overrides the `Insert` method, which emits a PubSub messages whenever `Insert` is
called, i.e. upon a successful `CheckTx` call.

To integrate into an application, simply set the mempool on `BaseApp`:

```go
import (
  // ...
  mempoolpubsub "github.com/numiadata/tools/mempool
)

func New(...) *App {
  // ...
  app.SetMempool(mempoolpubsub.NewPubSubMempool(logger, "<nodeID>", "<projectID>", "<topic>", false))
}
```

## PubSub Messages

The following PubSub messages are emitted:

```json
{
  "data": bytes(<MSG-STRING>),
  "attributes": {
    "message_type": "mempool_tx_msg",
    "chain_id": "<CHAIN-ID>",
    "tx_hash": "<TX-HASH>",
    "timestamp": "<TX-INSERT-TIMESTAMP-NANO>",
    "node_id": "<NODE-ID/MONIKER>",
    "msg_signer": "<TX-MSG-SIGNER>",
    "tx_msg_type": "<TX-MSG-TYPE>",
  }
}
```
