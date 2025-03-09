package entity

type Transfer struct {
	ID       uint64
	Sender   string
	Receiver string
	Amount   int64
	Made_At  string
}
