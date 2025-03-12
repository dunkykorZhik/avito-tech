package service

import (
	"context"
	"time"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type InfoOutput struct {
	Balance         int64             `json:"coins"`
	InventoryOutput []InventoryOutput `json:"inventory"`
	CoinHistory     CoinHistory       `json:"coinHistory"`
}
type InventoryOutput struct {
	Item_name string `json:"type"`
	Quantity  int64  `json:"quantity"`
}
type CoinHistory struct {
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
	ParseToken(tokenString string) (string, error)
}
type Transfer interface {
	CreateTransfer(ctx context.Context, transfer entity.Transfer) error
}
type History interface {
	GetHistory(ctx context.Context, username string) (*InfoOutput, error)
}
type Inventory interface {
	BuyItem(ctx context.Context, username, item_name string) error
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
		Inventory: NewInventoryService(deps.Repo.User, deps.Repo.Inventory, deps.Repo.Merch),
		History:   NewHistoryService(deps.Repo.User, deps.Repo.Transfer, deps.Repo.Inventory),
	}
}
