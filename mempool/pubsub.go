package mempool

const (
	// CredsEnvVar defines the environment variable that must be present when a using
	// Google Cloud PubSub client.
	CredsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

	// PubSub message attribute keys
	AttrKeyMsgType   = "message_type"
	AttrKeyChainID   = "chain_id"
	AttrKeyTxHash    = "tx_hash"
	AttrKeyTimestamp = "timestamp"
	AttrKeyNodeID    = "node_id"
	AttrKeyMsgSigner = "msg_signer"
	AttrKeyTxMsgType = "tx_msg_type"

	// MsgTypeCheckTxMsg defines the attribute value for 'AttrKeyMsgType'
	MsgTypeCheckTxMsg = "mempool_tx_msg"
)
