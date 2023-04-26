package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"strconv"
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

func (d *DebtPostgres) GetAllDebts(debtorId *int,
	lenderId *int,
	sortBy []string) ([]data.Debt, error) {

	query := fmt.Sprintf("SELECT * FROM %s ", debtsTable)
	if debtorId != nil || lenderId != nil {
		query += "WHERE "
	}
	var params []interface{}

	if debtorId != nil {
		query += " debtor_id = $" + strconv.Itoa(len(params)+1)
		params = append(params, debtorId)
	}

	if lenderId != nil {
		if len(params) > 0 {
			query += " AND "
		}
		query += " lender_id = $" + strconv.Itoa(len(params)+1)
		params = append(params, lenderId)
	}

	if len(sortBy) != 0 {
		query += fmt.Sprintf(" ORDER BY ")
		for index, value := range sortBy {
			if index > 0 {
				query += ", "
			}
			if value[0] == '-' {
				query += value[1:] + " DESC"
			} else {
				query += value
			}
		}
	}

	rows, err := d.db.Queryx(query, params...)
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
	query := fmt.Sprintf(
		"UPDATE %s SET status=$2 WHERE id=$1 ",
		debtsTable)

	_, err := d.db.Exec(query, id, status)
	return err
}

func (d *DebtPostgres) DeleteDebt(debtId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", debtsTable)
	_, err := d.db.Exec(query, debtId)
	return err

}
