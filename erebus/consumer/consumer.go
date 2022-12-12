package consumer

// Supported message types
const (
	MessageTypeData MessageType = iota
	MessageTypeMeta
)

type (
	MessageType uint32

	// Message defines a structure containing the type and contents of a state
	// streamed datum to be proxied to a consumer.
	Message struct {
		Type MessageType
		Body interface{}
	}

	// Consumer defines the interface for a consumer that takes messages and pushes
	// them to a sink.
	Consumer interface {
		Push(Message)
	}
)
