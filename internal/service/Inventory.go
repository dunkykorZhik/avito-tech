package service

import (
	"context"
	"fmt"

	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type InventoryService struct {
	userRepo      repo.User
	inventoryRepo repo.Inventory
	merchRepo     repo.Merch
}

func NewInventoryService(userRepo repo.User, inventoryRepo repo.Inventory) *InventoryService {
	return &InventoryService{
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) BuyItem(ctx context.Context, id int64, item_name string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	err := s.userRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		if err := s.buyItemProcess(ctxTx, id, item_name); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *InventoryService) buyItemProcess(ctxTx context.Context, id int64, item_name string) error {
	user, err := s.userRepo.GetUserById(ctxTx, id)
	if err != nil {
		return nil
	}
	merch, err := s.merchRepo.GetMerch(ctxTx, item_name)
	if err != nil {
		return err
	}
	if user.Balance < merch.Cost {
		return fmt.Errorf("not enough balance")
	}
	if err = s.userRepo.Withdraw(ctxTx, merch.Cost, id); err != nil {
		return err
	}
	if err = s.inventoryRepo.AddItem(ctxTx, id, item_name); err != nil {
		return err
	}
	return nil

}
