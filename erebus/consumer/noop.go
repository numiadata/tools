package consumer

var _ Consumer = (*NoOpConsumer)(nil)

// NoOpConsumer defines a consumer that drops all messages.
type NoOpConsumer struct{}

func (_ NoOpConsumer) Push(_ Message) {}
