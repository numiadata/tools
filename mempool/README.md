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

### Mempool

TODO...
