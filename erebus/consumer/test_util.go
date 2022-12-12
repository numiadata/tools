package consumer

import (
	"sync"
)

type TestConsumer struct {
	mu       sync.RWMutex
	dataMsgs []Message
	metaMsgs []Message
}

func NewTestConsumer() *TestConsumer {
	return &TestConsumer{
		dataMsgs: make([]Message, 0),
		metaMsgs: make([]Message, 0),
	}
}

func (tc *TestConsumer) Push(m Message) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	switch m.Type {
	case MessageTypeData:
		tc.dataMsgs = append(tc.dataMsgs, m)

	case MessageTypeMeta:
		tc.metaMsgs = append(tc.metaMsgs, m)
	}
}

func (tc *TestConsumer) DataMessages() []Message {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	msgs := make([]Message, len(tc.dataMsgs))
	copy(msgs, tc.dataMsgs)
	return msgs
}

func (tc *TestConsumer) MetaMessages() []Message {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	msgs := make([]Message, len(tc.metaMsgs))
	copy(msgs, tc.metaMsgs)
	return msgs
}
