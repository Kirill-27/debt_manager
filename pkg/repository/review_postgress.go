package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
	"strconv"
)

type ReviewPostgres struct {
	db *sqlx.DB
}

func NewReviewPostgres(db *sqlx.DB) *ReviewPostgres {
	return &ReviewPostgres{db: db}
}

func (r *ReviewPostgres) GetReviewById(id int) (*data.Review, error) {
	var review data.Review
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", reviewsTable)
	err := r.db.Get(&review, query, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &review, err
}

func (r *ReviewPostgres) UpdateReview(review data.Review) error {
	query := fmt.Sprintf("UPDATE %s SET comment=$2, rate=$3,  WHERE id=$1 ", reviewsTable)

	_, err := r.db.Exec(query, review.Comment, review.Rate)
	return err
}

func (r *ReviewPostgres) CreateReview(review data.Review) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (reviewer_id, lender_id, comment, rate)"+
		" values ($1, $2, $3, $4, $5) RETURNING id", reviewsTable)

	row := r.db.QueryRow(query,
		review.ReviewerId,
		review.LenderId,
		review.Comment,
		review.Rate)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ReviewPostgres) GetAllReviews(reviewerId *int, lenderId *int, sortBy []string) ([]data.Review, error) {
	query := fmt.Sprintf("SELECT * FROM %s ", currentDebtsTable)
	if reviewerId != nil || lenderId != nil {
		query += "WHERE "
	}
	var params []interface{}

	if reviewerId != nil {
		query += " reviewer_id = $" + strconv.Itoa(len(params)+1)
		params = append(params, reviewerId)
	}

	if lenderId != nil {
		if len(params) > 0 {
			query += " AND "
		}
		query += " lender_id = $" + strconv.Itoa(len(params)+1)
		params = append(params, lenderId)
	}

	if len(sortBy) != 0 {
		query += helpers.ParseSortBy(sortBy)
	}

	rows, err := r.db.Queryx(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []data.Review
	for rows.Next() {
		var review data.Review
		if err := rows.StructScan(&review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
