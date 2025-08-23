package service

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type InventoryService struct {
	userRepo      repo.User
	inventoryRepo repo.Inventory
}

func NewInventoryService(userRepo repo.User, inventoryRepo repo.Inventory) *InventoryService {
	return &InventoryService{
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) BuyItem(ctx context.Context, userID int64, item_name string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	err := s.userRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		if err := s.buyItemProcess(ctxTx, userID, item_name); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *InventoryService) buyItemProcess(ctxTx context.Context, userID int64, item_name string) error {
	itemID, err := s.userRepo.BuyItem(ctxTx, item_name, userID)
	if err != nil {
		return err
	}
	if err = s.inventoryRepo.AddItem(ctxTx, userID, itemID); err != nil {
		return err
	}
	return nil

}
