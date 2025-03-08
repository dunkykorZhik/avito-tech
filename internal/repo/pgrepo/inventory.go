package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
)

type inventoryRepo struct {
	db db.Database
}

func NewInventoryRepo(db db.Database) *inventoryRepo {
	return &inventoryRepo{db: db}

}

func (r *inventoryRepo) AddItem(ctxTx context.Context, user_id int64, item_name string) error {
	q, err := r.checkItemQuantity(ctxTx, user_id, item_name)
	if err != nil {
		return err
	}
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := ``
	if q > 0 {
		query = `UPDATE inventory SET quantity=quantity+1 WHERE user_id=$1 AND item_name=$2`
	} else {
		query = `INSERT INTO inventory (user_id, item_name, 1) VALUES ($1, $2)`
	}
	_, err = tx.ExecContext(ctxTx, query, user_id, item_name)
	if err != nil {
		return err
	}
	return nil

}
func (r *inventoryRepo) GetInventory(ctx context.Context, id int64) ([]entity.Inventory, error) {

	query := `SELECT id, user_id, item_name, quantity FROM inventory WHERE user_id = $1`
	inventories := []entity.Inventory{}
	rows, err := r.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i entity.Inventory
		if err := rows.Scan(&i.ID, &i.UserId, &i.ItemName, &i.Quantity); err != nil {
			return nil, err
		}
		inventories = append(inventories, i)

	}

	return inventories, nil

}
func (r *inventoryRepo) checkItemQuantity(ctxTx context.Context, user_id int64, item_name string) (int64, error) {
	query := `SELECT quantity FROM inventory WHERE user_id=$1 AND item_name=$2`
	var q int64
	if err := r.db.GetDb().QueryRowContext(ctxTx, query, user_id, item_name).Scan(&q); err != nil {
		return 0, err
	}
	return q, nil

}
