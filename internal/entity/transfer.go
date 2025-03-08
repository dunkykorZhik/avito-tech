package entity

type Transfer struct {
	ID         uint64
	SenderID   int64
	ReceiverID int64
	Amount     int64
	Made_At    string
}
