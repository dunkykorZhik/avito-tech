package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
)

type UserRepo struct {
	db db.Database
}

func NewUserRepo(db db.Database) *UserRepo {
	return &UserRepo{db: db}
}

type key string

const TxCtxKey key = "TxCtxKey"

// TODO hash the password
func (r *UserRepo) GetUserById(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, username, password, balance FROM users WHERE id = $1`
	var user entity.User
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func (r *UserRepo) Deposit(ctxTx context.Context, amount, id int64) error {
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `UPDATE users SET balance = balance+$1 WHERE id=%2`
	_, err := tx.ExecContext(ctxTx, query, amount, id)
	if err != nil {
		return err
	}
	return nil

}
func (r *UserRepo) Withdraw(ctxTx context.Context, amount, id int64) error {

	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `UPDATE users SET balance = balance-$1 WHERE id=%2`
	_, err := tx.ExecContext(ctxTx, query, amount, id)
	if err != nil {
		return err
	}
	return nil

}

func (r *UserRepo) WithTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	tx, err := r.db.GetDb().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting tx : %v", err)
	}
	ctxTx := context.WithValue(ctx, TxCtxKey, tx)
	err = fn(ctxTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling tx back : %v, which is happening cause %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()

}
