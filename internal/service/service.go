package service

import (
	"context"
	"time"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type InfoResponse struct {
	Balance         int64             `json:"coins"`
	InventoryOutput []InventoryOutput `json:"inventory"`
	TransferHistory TransferHistory   `json:"coinHistory"`
}
type InventoryOutput struct {
	Item_name string `json:"type"`
	Quantity  int64  `json:"quantity"`
}
type TransferHistory struct {
	SentOutput     []SentOutput     `json:"sent"`
	ReceivedOutput []ReceivedOutput `json:"received"`
}
type SentOutput struct {
	User_name string `json:"toUser"`
	Amount    int64  `json:"Amount"`
}
type ReceivedOutput struct {
	User_name string `json:"fromUser"`
	Amount    int64  `json:"Amount"`
}

const ctxTimeout = time.Second * 5

type User interface {
	GenerateToken(ctx context.Context, username string, password string) (string, error)
	ParseToken(tokenString string) (int64, string, error)
}
type Transfer interface {
	CreateTransfer(ctx context.Context, transfer entity.Transfer) error
}
type History interface {
	GetHistory(ctx context.Context, userID int64) (*InfoResponse, error)
}
type Inventory interface {
	BuyItem(ctx context.Context, userID int64, item_name string) error
}

type Service struct {
	User      User
	Transfer  Transfer
	Inventory Inventory
	History   History
}

type ServiceDependencies struct {
	Repo     *repo.Repositories
	SignKey  string
	TokenTTL time.Duration
}

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		User:      NewUserService(deps.Repo.User, deps.SignKey, deps.TokenTTL),
		Transfer:  NewTransferService(deps.Repo.User, deps.Repo.Transfer),
		Inventory: NewInventoryService(deps.Repo.User, deps.Repo.Inventory),
		History:   NewHistoryService(deps.Repo.User, deps.Repo.Transfer, deps.Repo.Inventory),
	}
}
