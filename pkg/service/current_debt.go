package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

type CurrentDebtService struct {
	repo repository.CurrentDebt
}

func NewCurrentDebtService(repo repository.CurrentDebt) *CurrentDebtService {
	return &CurrentDebtService{repo: repo}
}

func (s *CurrentDebtService) CreateCurrentDebt(debt data.CurrentDebts) (int, error) {
	return s.repo.CreateCurrentDebt(debt)
}

func (s *CurrentDebtService) GetAllCurrentDebts(debtorId *int, lenderId *int, sortBy []string) (
	[]data.CurrentDebts, error) {
	return s.repo.GetAllCurrentDebts(debtorId, lenderId, sortBy)
}

func (s *CurrentDebtService) UpdateAmount(id int, amount int) error {
	return s.repo.UpdateAmount(id, amount)
}

func (s *CurrentDebtService) DeleteCurrentDebt(debtId int) error {
	return s.repo.DeleteCurrentDebt(debtId)
}
