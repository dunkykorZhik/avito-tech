package entity

type Transfer struct {
	ID           uint64
	SenderID     int64
	ReceiverID   int64
	SenderName   string
	ReceiverName string
	Amount       int64
	Made_At      string
}
