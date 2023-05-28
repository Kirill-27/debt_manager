package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
)

type StripePaymentPostgres struct {
	db *sqlx.DB
}

func NewStripePaymentPostgres(db *sqlx.DB) *StripePaymentPostgres {
	return &StripePaymentPostgres{db: db}
}

func (r *StripePaymentPostgres) CreateStripePayment(stripePayment data.StripePayment) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (payment_id, user_id, status)"+
		" values ($1, $2, $3) RETURNING payment_id", stripePaymentsTable)

	row := r.db.QueryRow(query,
		stripePayment.PaymentId,
		stripePayment.UserId,
		stripePayment.Status)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *StripePaymentPostgres) UpdateStripePaymentStatus(stripePaymentsId string, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status=$2 WHERE payment_id=$1 ", stripePaymentsTable)
	_, err := r.db.Exec(query, stripePaymentsId, status)
	return err
}

func (r *StripePaymentPostgres) GetStripePaymentByPaymentId(paymentId string) (*data.StripePayment, error) {
	var stripePayment data.StripePayment
	query := fmt.Sprintf("SELECT * FROM %s WHERE payment_id=$1", stripePaymentsTable)
	err := r.db.Get(&stripePayment, query, paymentId)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &stripePayment, err
}

func (r *StripePaymentPostgres) GetAllStripePayments(status *int, sortBy []string) ([]data.StripePayment, error) {
	query := fmt.Sprintf("SELECT * FROM %s ", stripePaymentsTable)
	var params []interface{}

	if status != nil {
		query += "WHERE status=$1"
		params = append(params, status)
	}

	if len(sortBy) != 0 {
		query += helpers.ParseSortBy(sortBy)
	}

	rows, err := r.db.Queryx(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stripePayments []data.StripePayment
	for rows.Next() {
		var stripePayment data.StripePayment
		if err := rows.StructScan(&stripePayment); err != nil {
			return nil, err
		}
		stripePayments = append(stripePayments, stripePayment)
	}

	return stripePayments, nil
}
