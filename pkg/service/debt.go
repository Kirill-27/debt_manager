package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

type DebtService struct {
	repo repository.Debt
}

func NewDebtService(repo repository.Debt) *DebtService {
	return &DebtService{repo: repo}
}

func (s *DebtService) CreateDebt(debt data.Debt) (int, error) {
	return s.repo.CreateDebt(debt)
}

func (s *DebtService) GetAllDebts(debtorId *int, lenderId *int, statuses string, sortBy []string) ([]data.Debt, error) {
	return s.repo.GetAllDebts(debtorId, lenderId, statuses, sortBy)
}
func (s *DebtService) GetDebtById(debtId int) (*data.Debt, error) {
	return s.repo.GetDebtById(debtId)
}
func (s *DebtService) UpdateStatus(id int, status int) error {
	return s.repo.UpdateStatus(id, status)
}
func (s *DebtService) DeleteDebt(debtId int) error {
	return s.repo.DeleteDebt(debtId)
}

func (s *DebtService) CloseAllDebts(debtorId int, lenderId int) error {
	return s.repo.CloseAllDebts(debtorId, lenderId)
}

func (s *DebtService) SelectTop3Lenders(debtorId int) ([]int, error) {
	return s.repo.SelectTop3Lenders(debtorId)
}

func (s *DebtService) SelectTop3Debtors(lenderId int) ([]int, error) {
	return s.repo.SelectTop3Debtors(lenderId)
}
