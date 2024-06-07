package repository

import (
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"database/sql"
	"log"
	"math"
	"time"
)

type expensesRepo struct {
	db *sql.DB
}

func (a *expensesRepo) InsertExpenses(expenses model.Expenses) (model.Expenses, error) {
	stmt, err := a.db.Prepare("INSERT INTO expenses(date,amount,transaction_type,balance, description, user_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return model.Expenses{}, err
	}
	defer stmt.Close()

	var expensesId string
	err = stmt.QueryRow(expenses.Date, expenses.Amount, expenses.Transaction_type, expenses.Balance, expenses.Description, expenses.User_Id).Scan(&expensesId)
	if err != nil {
		return model.Expenses{}, err
	}

	expenses.Id = expensesId
	expenses.Date = time.Now()

	expenses.CreatedAt = time.Now()
	expenses.UpdatedAt = time.Now()
	return expenses, nil
}

func (a *expensesRepo) FindAllExpenses(page int, size int) ([]model.Expenses, dto.Paging, error) {
	var listData []model.Expenses
	var row *sql.Rows

	// rumus pagination
	offset := (page - 1) * size

	var err error
	row, err = a.db.Query("select * from expenses limit $1 offset $2", size, offset)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	totalRows := 0
	err = a.db.QueryRow("SELECT COUNT(*) FROM expenses").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	for row.Next() {
		var expenses model.Expenses
		err := row.Scan(&expenses.Id, &expenses.Date, &expenses.Amount, &expenses.Transaction_type, &expenses.Balance, &expenses.Description, &expenses.CreatedAt,
			&expenses.UpdatedAt, &expenses.User_Id)
		if err != nil {
			log.Println(err.Error())
		}
		listData = append(listData, expenses)
	}
	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return listData, paging, nil
}

func (a *expensesRepo) FindByIdExpenses(id string) (model.Expenses, error) {
	var expenses model.Expenses
	err := a.db.QueryRow("select * from expenses where id=$1", id).
		Scan(&expenses.Id, &expenses.Date, &expenses.Amount, &expenses.Transaction_type,
			&expenses.Balance, &expenses.Description, &expenses.CreatedAt, &expenses.UpdatedAt, &expenses.User_Id)
	if err != nil {
		return model.Expenses{}, err
	}
	return expenses, nil
}

func (a *expensesRepo) FindByTypeExpenses(typ string) ([]model.Expenses, error) {
	rows, err := a.db.Query("SELECT * from expenses where transaction_type=$1", typ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expenses
	for rows.Next() {
		var exp model.Expenses
		err := rows.Scan(&exp.Id, &exp.Date, &exp.Amount, &exp.Transaction_type,
			&exp.Balance, &exp.Description, &exp.CreatedAt, &exp.UpdatedAt, &exp.User_Id)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, exp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}
func (a *expensesRepo) UpdatedExpenses(expenses model.Expenses) (model.Expenses, error) {
	stmt, err := a.db.Prepare(`
        UPDATE expenses 
        SET amount = $1, transaction_type = $2, balance = $3, description = $4, updated_at = $5, user_id = $6
        WHERE id = $7 
        RETURNING id, date, amount, transaction_type, balance, description, created_at, updated_at, user_id
    `)
	if err != nil {
		return model.Expenses{}, err
	}
	defer stmt.Close()

	var updatedExpense model.Expenses
	err = stmt.QueryRow(expenses.Amount, expenses.Transaction_type, expenses.Balance, expenses.Description, time.Now(), expenses.User_Id, expenses.Id).
		Scan(&updatedExpense.Id, &updatedExpense.Date, &updatedExpense.Amount, &updatedExpense.Transaction_type, &updatedExpense.Balance, &updatedExpense.Description, &updatedExpense.CreatedAt, &updatedExpense.UpdatedAt, &updatedExpense.User_Id)
	if err != nil {
		return model.Expenses{}, err
	}

	return updatedExpense, nil
}

type ExpensesRepo interface {
	FindAllExpenses(page int, size int) ([]model.Expenses, dto.Paging, error)
	FindByIdExpenses(id string) (model.Expenses, error)
	FindByTypeExpenses(id string) ([]model.Expenses, error)
	InsertExpenses(expenses model.Expenses) (model.Expenses, error)
	UpdatedExpenses(expenses model.Expenses) (model.Expenses, error)
}

func NewExpensesRepo(database *sql.DB) ExpensesRepo {
	return &expensesRepo{db: database}
}
