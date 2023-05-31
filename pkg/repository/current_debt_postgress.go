package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
	"strconv"
)

type CurrentDebtPostgres struct {
	db *sqlx.DB
}

func NewCurrentDebtPostgres(db *sqlx.DB) *CurrentDebtPostgres {
	return &CurrentDebtPostgres{db: db}
}

func (d *CurrentDebtPostgres) CreateCurrentDebt(debt data.CurrentDebts) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (debtor_id, lender_id, amount)"+
		" values ($1, $2, $3) RETURNING id", currentDebtsTable)

	row := d.db.QueryRow(query,
		debt.DebtorID,
		debt.LenderId,
		debt.Amount)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (d *CurrentDebtPostgres) GetAllCurrentDebts(debtorId *int, lenderId *int, sortBy []string) (
	[]data.CurrentDebts, error) {

	query := fmt.Sprintf("SELECT * FROM %s ", currentDebtsTable)
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
		query += helpers.ParseSortBy(sortBy)
	}

	rows, err := d.db.Queryx(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debts []data.CurrentDebts
	for rows.Next() {
		var debt data.CurrentDebts
		if err := rows.StructScan(&debt); err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}

	return debts, nil
}

func (d *CurrentDebtPostgres) UpdateAmount(id int, amount int) error {
	query := fmt.Sprintf(
		"UPDATE %s SET amount=$2 WHERE id=$1 ",
		currentDebtsTable)

	_, err := d.db.Exec(query, id, amount)
	return err
}

func (d *CurrentDebtPostgres) DeleteCurrentDebt(debtId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", currentDebtsTable)
	_, err := d.db.Exec(query, debtId)
	return err
}
