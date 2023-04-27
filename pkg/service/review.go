package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

type ReviewService struct {
	repo repository.Review
}

func NewReviewService(repo repository.Review) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) GetReviewById(id int) (*data.Review, error) {
	return s.repo.GetReviewById(id)
}

func (s *ReviewService) GetAllReviews(reviewerId *int, lenderId *int, sortBy []string) ([]data.Review, error) {
	return s.repo.GetAllReviews(reviewerId, lenderId, sortBy)
}

func (s *ReviewService) UpdateReview(review data.Review) error {
	return s.repo.UpdateReview(review)
}

func (s *ReviewService) CreateReview(review data.Review) (int, error) {
	return s.repo.CreateReview(review)
}
