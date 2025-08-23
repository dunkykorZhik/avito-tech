package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/app/db"
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
	_, err := tx.ExecContext(ctxTx, query, transfer.SenderID, transfer.ReceiverID, transfer.Amount)
	if err != nil {
		return err
	}
	return nil

}

func (t *transferRepo) GetSentHistory(ctx context.Context, id int64) ([]entity.Transfer, error) {
	query := `SELECT  t.id, r.username AS receiver_username, t.amount, t.made_at FROM transfers t JOIN users r ON t.receiver = r.id WHERE t.sender = $1`
	transfers := []entity.Transfer{}
	rows, err := t.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t entity.Transfer
		if err := rows.Scan(&t.ID, &t.ReceiverName, &t.Amount, &t.Made_At); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)

	}

	return transfers, nil

}
func (t *transferRepo) GetReceivedHistory(ctx context.Context, id int64) ([]entity.Transfer, error) {
	query := `SELECT  t.id, s.username AS sender_username, t.amount, t.made_at FROM transfers t JOIN users s ON t.sender = s.id WHERE t.receiver = $1`
	transfers := []entity.Transfer{}
	rows, err := t.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t entity.Transfer
		if err := rows.Scan(&t.ID, &t.SenderName, &t.Amount, &t.Made_At); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)

	}

	return transfers, nil

}
