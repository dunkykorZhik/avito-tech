package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
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
func (r *UserRepo) GetUserByName(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, password, balance FROM users WHERE username = $1`
	var user entity.User
	err := r.db.GetDb().QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoUser
		}
		return nil, err
	}
	return &user, nil

}

func (r *UserRepo) CreateUser(ctx context.Context, username, password string) error {

	query := `INSERT INTO users (username, password, balance) VALUES ( $1, $2, 1000)`
	_, err := r.db.GetDb().ExecContext(ctx, query, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Deposit(ctxTx context.Context, amount int64, username string) error {
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `UPDATE users SET balance = balance+$1 WHERE username=$2`
	res, err := tx.ExecContext(ctxTx, query, amount, username)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return errs.ErrNoUser
	}
	return nil

}
func (r *UserRepo) Withdraw(ctxTx context.Context, amount int64, username string) error {

	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `UPDATE users SET balance = balance-$1 WHERE username=$2`
	res, err := tx.ExecContext(ctxTx, query, amount, username)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return errs.ErrNoUser
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
