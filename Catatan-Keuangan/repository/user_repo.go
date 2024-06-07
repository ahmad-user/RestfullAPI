package repository

import (
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"database/sql"
	"log"
	"math"
)

type userRepo struct {
	db *sql.DB
}

func (a *userRepo) FindAll(page int, size int) ([]model.Users, dto.Paging, error) {
	var listData []model.Users
	var row *sql.Rows
	offset := (page - 1) * size
	var err error
	row, err = a.db.Query("select * from users limit $1 offset $2", size, offset)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	totalRows := 0
	err = a.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	for row.Next() {
		var user model.Users

		err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreateAt, &user.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
		}
		listData = append(listData, user)
	}
	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return listData, paging, nil
}

func (a *userRepo) FindById(id string) (model.Users, error) {
	var user model.Users
	err := a.db.QueryRow("select * from users where id=$1", id).
		Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreateAt, &user.UpdatedAt)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

func (a *userRepo) FindByUsername(username string) (model.Users, error) {
	var user model.Users
	err := a.db.QueryRow("select * from users where username=$1", username).
		Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreateAt, &user.UpdatedAt)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

type UserRepo interface {
	FindAll(page int, size int) ([]model.Users, dto.Paging, error)
	FindById(id string) (model.Users, error)
	FindByUsername(username string) (model.Users, error)
}

func NewUserRepo(database *sql.DB) UserRepo {
	return &userRepo{db: database}
}
