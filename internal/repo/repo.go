package repo

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo/pgrepo"
)

type User interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUserByName(ctx context.Context, username string) (*entity.User, error)
	//GetUserById(ctx context.Context, id int64) (*entity.User, error)

	Deposit(ctxTx context.Context, amount int64, username string) error
	Withdraw(ctxTx context.Context, amount int64, username string) error
	WithTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Transfer interface {
	CreateTransfer(ctxTx context.Context, transaction entity.Transfer) error
	GetSentHistory(ctx context.Context, username string) ([]entity.Transfer, error)
	GetReceivedHistory(ctx context.Context, username string) ([]entity.Transfer, error)

	// check balance->
	/*

	 */
}

type Inventory interface {
	AddItem(ctxTx context.Context, username string, item_name string) error
	GetInventory(ctx context.Context, username string) ([]entity.Inventory, error)
}
type Merch interface {
	GetMerch(ctx context.Context, item_name string) (*entity.Merch, error)
}
type Repositories struct {
	User
	Transfer
	Inventory
	Merch
}

// MIGHTDO: change to factory or better injection. repo do not know which db
func NewRepositories(db db.Database) *Repositories {
	return &Repositories{
		User:      pgrepo.NewUserRepo(db),
		Transfer:  pgrepo.NewTransferRepo(db),
		Inventory: pgrepo.NewInventoryRepo(db),
		Merch:     pgrepo.NewMerchRepo(db),
	}
}
