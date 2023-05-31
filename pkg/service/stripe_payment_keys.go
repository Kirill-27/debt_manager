package service

import (
	"github.com/kirill-27/debt_manager/pkg/repository"
	"time"
)

type StripePaymentKeysService struct {
	repo repository.StripePaymentKeys
}

func NewStripePaymentKeyService(repo repository.StripePaymentKeys) *StripePaymentKeysService {
	return &StripePaymentKeysService{repo: repo}
}

func (s *StripePaymentKeysService) GetLastHandled() (time.Time, error) {
	return s.repo.GetLastHandled()
}
func (s *StripePaymentKeysService) SetLastHandled(lastHandledTime int64) error {
	return s.repo.SetLastHandled(lastHandledTime)
}
