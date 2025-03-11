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

func (s *HistoryService) GetHistory(ctx context.Context, username string) (*InfoOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	user, err := s.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	sent, err := s.transferRepo.GetSentHistory(ctx, username)
	if err != nil {
		return nil, err
	}
	received, err := s.transferRepo.GetReceivedHistory(ctx, username)
	if err != nil {
		return nil, err
	}
	inventory, err := s.inventoryRepo.GetInventory(ctx, username)
	if err != nil {
		return nil, err
	}
	return createOutput(user.Balance, sent, received, inventory), nil

}

func createOutput(balance int64, sent, received []entity.Transfer, inventory []entity.Inventory) *InfoOutput {

	s := make([]SentOutput, 0, len(sent))
	for _, transfer := range sent {
		s = append(s, SentOutput{
			transfer.Receiver,
			transfer.Amount,
		})

	}
	r := make([]ReceivedOutput, 0, len(received))
	for _, transfer := range received {
		r = append(r, ReceivedOutput{
			transfer.Sender,
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

	return &InfoOutput{
		Balance:         balance,
		CoinHistory:     CoinHistory{s, r},
		InventoryOutput: i,
	}

}
