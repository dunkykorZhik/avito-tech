package repo

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo/pgrepo"
)

type User interface {
	CreateUser(ctx context.Context, username string, password string) (int64, error)
	GetUserByName(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)

	Deposit(ctxTx context.Context, amount int64, username string) (int64, error)
	Withdraw(ctxTx context.Context, amount int64, id int64) error
	BuyItem(ctxTx context.Context, item_name string, id int64) (int64, error)
	WithTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Transfer interface {
	CreateTransfer(ctxTx context.Context, transaction entity.Transfer) error
	GetSentHistory(ctx context.Context, id int64) ([]entity.Transfer, error)
	GetReceivedHistory(ctx context.Context, id int64) ([]entity.Transfer, error)

	// check balance->
	/*

	 */
}

type Inventory interface {
	AddItem(ctxTx context.Context, userID, itemID int64) error
	GetInventory(ctx context.Context, id int64) ([]entity.Inventory, error)
}

type Repositories struct {
	User
	Transfer
	Inventory
}

// MIGHTDO: change to factory or better injection. repo do not know which db
func NewRepositories(db db.Database) *Repositories {
	return &Repositories{
		User:      pgrepo.NewUserRepo(db),
		Transfer:  pgrepo.NewTransferRepo(db),
		Inventory: pgrepo.NewInventoryRepo(db),
	}
}
