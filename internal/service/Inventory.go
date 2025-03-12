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

func NewInventoryService(userRepo repo.User, inventoryRepo repo.Inventory, merchRepo repo.Merch) *InventoryService {
	return &InventoryService{
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
		merchRepo:     merchRepo,
	}
}

func (s *InventoryService) BuyItem(ctx context.Context, username, item_name string) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	err := s.userRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		if err := s.buyItemProcess(ctxTx, username, item_name); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *InventoryService) buyItemProcess(ctxTx context.Context, username, item_name string) error {
	user, err := s.userRepo.GetUserByName(ctxTx, username)
	if err != nil {
		return err
	}
	merch, err := s.merchRepo.GetMerch(ctxTx, item_name)
	if err != nil {
		return err
	}
	if user.Balance < merch.Cost {
		return fmt.Errorf("not enough balance")
	}
	if err = s.userRepo.Withdraw(ctxTx, merch.Cost, username); err != nil {
		return err
	}
	if err = s.inventoryRepo.AddItem(ctxTx, username, item_name); err != nil {
		return err
	}
	return nil

}
