package usecase

import (
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"RestFullAPI/repository"
)

type userUseCase struct {
	repo repository.UserRepo
}

func (a *userUseCase) FindAll(page int, size int) ([]model.Users, dto.Paging, error) {
	return a.repo.FindAll(page, size)
}

func (a *userUseCase) FindById(id string) (model.Users, error) {
	return a.repo.FindById(id)
}

func (a *userUseCase) FindByUsername(username string) (model.Users, error) {
	return a.repo.FindByUsername(username)
}

type UserUseCase interface {
	// authors
	FindAll(page int, size int) ([]model.Users, dto.Paging, error)
	FindById(id string) (model.Users, error)
	FindByUsername(username string) (model.Users, error)
	// task
}

func MenuAuthorUseCase(repo repository.UserRepo) UserUseCase {
	return &userUseCase{repo: repo}
}
