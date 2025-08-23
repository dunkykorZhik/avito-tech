package service

import (
	"context"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
)

type TransferService struct {
	userRepo     repo.User
	transferRepo repo.Transfer
}

func NewTransferService(userRepo repo.User, transferRepo repo.Transfer) *TransferService {
	return &TransferService{userRepo: userRepo, transferRepo: transferRepo}
}

func (s *TransferService) CreateTransfer(ctx context.Context, transfer entity.Transfer) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	err := s.userRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		err := s.createTransferProcess(ctxTx, transfer)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *TransferService) createTransferProcess(ctx context.Context, transfer entity.Transfer) error {

	if err := s.userRepo.Withdraw(ctx, transfer.Amount, transfer.SenderID); err != nil {
		return err
	}
	var err error
	transfer.ReceiverID, err = s.userRepo.Deposit(ctx, transfer.Amount, transfer.ReceiverName)
	if err != nil {
		return err
	}
	if err := s.transferRepo.CreateTransfer(ctx, transfer); err != nil {
		return err
	}
	return nil
}
