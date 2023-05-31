package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

type StripePaymentService struct {
	repo repository.StripePayment
}

func NewStripePaymentService(repo repository.StripePayment) *StripePaymentService {
	return &StripePaymentService{repo: repo}
}

func (s *StripePaymentService) CreateStripePayment(stripePayment data.StripePayment) (string, error) {
	return s.repo.CreateStripePayment(stripePayment)
}

func (s *StripePaymentService) UpdateStripePaymentStatus(stripePaymentsId string, status int) error {
	return s.repo.UpdateStripePaymentStatus(stripePaymentsId, status)
}
func (s *StripePaymentService) GetStripePaymentByPaymentId(paymentId string) (*data.StripePayment, error) {
	return s.repo.GetStripePaymentByPaymentId(paymentId)
}

func (s *StripePaymentService) GetAllStripePayments(status *int, sortBy []string) ([]data.StripePayment, error) {
	return s.repo.GetAllStripePayments(status, sortBy)
}
