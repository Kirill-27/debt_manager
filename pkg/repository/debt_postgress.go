package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
)

type DebtPostgres struct {
	db *sqlx.DB
}

func NewDebtPostgres(db *sqlx.DB) *DebtPostgres {
	return &DebtPostgres{db: db}
}

func (d *DebtPostgres) CreateDebt(debt data.Debt) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (debtor_id, lender_id, amount, description)"+
		" values ($1, $2, $3, $4) RETURNING id", debtsTable)

	row := d.db.QueryRow(query,
		debt.DebtorId,
		debt.LenderId,
		debt.Amount,
		debt.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DebtPostgres) GetAllDebts(debtorIds string, lenderIds string, statuses string, sortBy []string) (
	[]data.Debt, error) {

	query := fmt.Sprintf("SELECT * FROM %s ", debtsTable)
	if debtorIds != "" || lenderIds != "" || statuses != "" {
		query += "WHERE "
	}
	params := 0

	if debtorIds != "" {
		query += " debtor_id IN (" + debtorIds + ")"
		params++
	}

	if lenderIds != "" {
		if params > 0 {
			query += " AND "
		}
		query += " lender_id IN (" + lenderIds + ")"
		params++
	}

	// it must be the last check for where statement
	if statuses != "" {
		if params > 0 {
			query += " AND "
		}
		query += " status IN (" + statuses + ")"
	}

	if len(sortBy) != 0 {
		query += helpers.ParseSortBy(sortBy)
	}

	rows, err := d.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debts []data.Debt
	for rows.Next() {
		var debt data.Debt
		if err := rows.StructScan(&debt); err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}

	return debts, nil
}

func (d *DebtPostgres) GetDebtById(debtId int) (*data.Debt, error) {
	var debt data.Debt
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", debtsTable)
	err := d.db.Get(&debt, query, debtId)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &debt, err
}
func (d *DebtPostgres) UpdateStatus(id int, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status=$2 WHERE id=$1 ", debtsTable)
	_, err := d.db.Exec(query, id, status)
	return err
}

func (d *DebtPostgres) DeleteDebt(debtId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", debtsTable)
	_, err := d.db.Exec(query, debtId)
	return err
}

func (d *DebtPostgres) CloseAllDebts(debtorId int, lenderId int) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1 WHERE status=$4"+
		" AND ((debtor_id=$2 AND lender_id=$3) OR  (debtor_id=$3 AND lender_id=$2))", debtsTable)
	_, err := d.db.Exec(query, data.DebtStatusClosed, debtorId, lenderId, data.DebtStatusActive)
	return err
}

func (d *DebtPostgres) SelectTop3Lenders(debtorId int) ([]int, error) {
	query := fmt.Sprintf("SELECT lender_id FROM %s WHERE debtor_id=$1 AND status in (2,3)"+
		" GROUP BY lender_id ORDER BY COUNT(*) DESC LIMIT 3", debtsTable)

	rows, err := d.db.Queryx(query, debtorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lenderIds []int
	for rows.Next() {
		var lenderId LenderIdStruct
		if err := rows.StructScan(&lenderId); err != nil {
			return nil, err
		}
		lenderIds = append(lenderIds, lenderId.LenderId)
	}

	return lenderIds, err
}

func (d *DebtPostgres) SelectTop3Debtors(lenderId int) ([]int, error) {
	query := fmt.Sprintf("SELECT debtor_id FROM %s WHERE lender_id=$1 AND status in (2,3)"+
		" GROUP BY debtor_id ORDER BY COUNT(*) DESC LIMIT 3", debtsTable)

	rows, err := d.db.Queryx(query, lenderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debtorIds []int
	for rows.Next() {
		var debtorId DebtorIdStruct
		if err := rows.StructScan(&debtorId); err != nil {
			return nil, err
		}
		debtorIds = append(debtorIds, debtorId.DebtorId)
	}

	return debtorIds, err
}

type LenderIdStruct struct {
	LenderId int `json:"lender_id" db:"lender_id"`
}

type DebtorIdStruct struct {
	DebtorId int `json:"debtor_id" db:"debtor_id"`
}
