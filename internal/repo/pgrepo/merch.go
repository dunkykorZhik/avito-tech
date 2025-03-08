package pgrepo

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/entity"
)

type merchRepo struct {
	db db.Database
}

func NewMerchRepo(db db.Database) *merchRepo {
	return &merchRepo{db: db}

}

func (r *merchRepo) GetMerch(ctx context.Context, item_name string) (*entity.Merch, error) {
	query := `SELECT id, item_name, cost FROM merch WHERE item_name = $1`
	var merch entity.Merch
	err := r.db.GetDb().QueryRowContext(ctx, query, item_name).Scan(&merch.Id, &merch.ItemName, &merch.Cost)
	if err != nil {
		return nil, err
	}
	return &merch, nil
}
