package service

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type HistoryService struct {
	userRepo      repo.User
	transferRepo  repo.Transfer
	inventoryRepo repo.Inventory
}

func NewHistoryService(userRepo repo.User, transferRepo repo.Transfer, inventoryRepo repo.Inventory) *HistoryService {
	return &HistoryService{
		userRepo:      userRepo,
		transferRepo:  transferRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *HistoryService) GetHistory(ctx context.Context, userID int64) (*InfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	sent, err := s.transferRepo.GetSentHistory(ctx, userID)
	if err != nil {
		return nil, err
	}
	received, err := s.transferRepo.GetReceivedHistory(ctx, userID)
	if err != nil {
		return nil, err
	}
	inventory, err := s.inventoryRepo.GetInventory(ctx, userID)
	if err != nil {
		return nil, err
	}
	return createOutput(user.Balance, sent, received, inventory), nil

}

func createOutput(balance int64, sent, received []entity.Transfer, inventory []entity.Inventory) *InfoResponse {

	s := make([]SentOutput, 0, len(sent))
	for _, transfer := range sent {
		s = append(s, SentOutput{
			transfer.ReceiverName,
			transfer.Amount,
		})

	}
	r := make([]ReceivedOutput, 0, len(received))
	for _, transfer := range received {
		r = append(r, ReceivedOutput{
			transfer.SenderName,
			transfer.Amount,
		})

	}

	i := make([]InventoryOutput, 0, len(inventory))
	for _, item := range inventory {
		i = append(i, InventoryOutput{
			Item_name: item.ItemName,
			Quantity:  item.Quantity,
		})
	}

	return &InfoResponse{
		Balance:         balance,
		TransferHistory: TransferHistory{s, r},
		InventoryOutput: i,
	}

}
