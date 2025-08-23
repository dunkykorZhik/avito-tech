package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
)

type inventoryRepo struct {
	db db.Database
}

func NewInventoryRepo(db db.Database) *inventoryRepo {
	return &inventoryRepo{db: db}

}

func (r *inventoryRepo) AddItem(ctxTx context.Context, userID, itemID int64) error {
	tx, ok := ctxTx.Value(TxCtxKey).(*sql.Tx)
	if !ok {
		return fmt.Errorf("no active transaction")
	}
	query := `INSERT INTO inventory (userID, itemID, quantity) VALUES ($1, $2, 1) ON CONFLICT (userID, itemID) DO UPDATE SET quantity = inventory.quantity + 1`

	_, err := tx.ExecContext(ctxTx, query, userID, itemID)
	if err != nil {
		return err
	}
	return nil

}
func (r *inventoryRepo) GetInventory(ctx context.Context, id int64) ([]entity.Inventory, error) {

	query := `SELECT m.item_name, i.quantity FROM inventory i JOIN merch m ON i.itemID = m.item_id WHERE i.userID = $1`
	inventories := []entity.Inventory{}
	rows, err := r.db.GetDb().QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i entity.Inventory
		if err := rows.Scan(&i.ID, &i.ItemName, &i.Quantity); err != nil {
			return nil, err
		}
		inventories = append(inventories, i)

	}

	return inventories, nil

}
