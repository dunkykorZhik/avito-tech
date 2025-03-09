package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
)

type transferRepo struct {
	db db.Database
}

func NewTransferRepo(db db.Database) *transferRepo {
	return &transferRepo{db: db}

}

func (t *transferRepo) CreateTransfer(ctxTx context.Context, transfer entity.Transfer) error {
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `INSERT INTO transfers (sender, receiver, amount) VALUES ($1,$2,$3)`
	_, err := tx.ExecContext(ctxTx, query, transfer.Sender, transfer.Receiver, transfer.Amount)
	if err != nil {
		return err
	}
	return nil

}

func (t *transferRepo) GetSentHistory(ctx context.Context, id int64) ([]entity.Transfer, error) {
	query := `SELECT id, sender, receiver, amount, made_at FROM transfers WHERE sender_id = $1`
	transfers := []entity.Transfer{}
	rows, err := t.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t entity.Transfer
		if err := rows.Scan(&t.ID, &t.Sender, &t.Receiver, &t.Amount, &t.Made_At); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)

	}

	return transfers, nil

}
func (t *transferRepo) GetReceivedHistory(ctx context.Context, id int64) ([]entity.Transfer, error) {
	query := `SELECT id, sender, receiver, amount, made_at FROM transfers WHERE receiver_id=$1`
	transfers := []entity.Transfer{}
	rows, err := t.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t entity.Transfer
		if err := rows.Scan(&t.ID, &t.Sender, &t.Receiver, &t.Amount, &t.Made_At); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)

	}

	return transfers, nil

}
