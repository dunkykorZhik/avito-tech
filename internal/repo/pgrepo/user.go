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

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, username, balance FROM users WHERE id = $1`
	var user entity.User
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoUser
		}
		return nil, err
	}
	return &user, nil

}

func (r *UserRepo) CreateUser(ctx context.Context, username, password string) (int64, error) {

	query := `INSERT INTO users (username, password, balance) VALUES ( $1, $2, 1000) RETURNING id`
	var id int64
	err := r.db.GetDb().QueryRowContext(ctx, query, username, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepo) Deposit(ctxTx context.Context, amount int64, username string) (int64, error) {
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return 0, fmt.Errorf("no active transaction")
	}
	var id int64
	query := `UPDATE users SET balance = balance+$1 WHERE username=$2 RETURNING id`
	err := tx.QueryRowContext(ctxTx, query, amount, username).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errs.ErrNoUser
	}
	return 0, nil

}
func (r *UserRepo) Withdraw(ctxTx context.Context, amount int64, id int64) error {

	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `UPDATE users SET balance = balance-$1 WHERE id=$2 AND balance>=$1`
	res, err := tx.ExecContext(ctxTx, query, amount, id)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return errs.ErrNotEnoughBalance
	}
	return nil

}

func (r *UserRepo) BuyItem(ctxTx context.Context, item_name string, id int64) (int64, error) {

	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return 0, fmt.Errorf("no active transaction")
	}

	var item_id, cost int64

	queryCost := `SELECT item_id, cost FROM merch WHERE item_name = $1`
	err := tx.QueryRowContext(ctxTx, queryCost, item_name).Scan(&item_id, &cost)
	if err == sql.ErrNoRows {
		return 0, errs.ErrNoItem
	}
	if err != nil {
		return 0, err
	}

	query := `UPDATE users SET balance = balance-$1 WHERE id=$2 AND balance>=$1`
	res, er := tx.ExecContext(ctxTx, query, cost, id)
	if er != nil {
		return 0, er
	}
	num, er := res.RowsAffected()
	if er != nil {
		return 0, er
	}
	if num == 0 {
		return 0, errs.ErrNotEnoughBalance
	}
	return id, nil

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
