package consumer

// Bigtable defines a consumer that pushes messages to Google Bigtable.
type Bigtable struct{}

func (b *Bigtable) Push(m Message) {
	panic("not implemented")
}
