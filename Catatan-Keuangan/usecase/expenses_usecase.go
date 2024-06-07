package usecase

import (
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"RestFullAPI/repository"
)

type expensesUseCase struct {
	expenses repository.ExpensesRepo
}

func (a *expensesUseCase) FindAllExpenses(page int, size int) ([]model.Expenses, dto.Paging, error) {
	return a.expenses.FindAllExpenses(page, size)
}

func (a *expensesUseCase) FindByIdExpenses(id string) (model.Expenses, error) {
	return a.expenses.FindByIdExpenses(id)
}
func (a *expensesUseCase) FindByTypeExpenses(typ string) ([]model.Expenses, error) {
	return a.expenses.FindByTypeExpenses(typ)
}
func (a *expensesUseCase) InsertExpenses(expenses model.Expenses) (model.Expenses, error) {
	return a.expenses.InsertExpenses(expenses)
}
func (a *expensesUseCase) UpdatedExpenses(expense model.Expenses) (model.Expenses, error) {
	return a.expenses.UpdatedExpenses(expense)
}

type ExpensesUseCase interface {
	FindAllExpenses(page int, size int) ([]model.Expenses, dto.Paging, error)
	FindByIdExpenses(id string) (model.Expenses, error)
	FindByTypeExpenses(id string) ([]model.Expenses, error)
	InsertExpenses(expenses model.Expenses) (model.Expenses, error)
	UpdatedExpenses(expense model.Expenses) (model.Expenses, error)
}

func MenuTasUserCase(repoExpenses repository.ExpensesRepo) ExpensesUseCase {
	return &expensesUseCase{expenses: repoExpenses}
}
