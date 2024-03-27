package user

import (
	"Encode/model/book"

	"gorm.io/gorm"
)

type Tbl_user struct {
	Nama     string
	Hp       string `gorm:"type:varchar(13);primaryKey"`
	Password string
	Book     []book.Tbl_book `gorm:"foreignKey:Pemilik;references:Hp"`
}

type UserModel struct {
	Connection *gorm.DB
}

func (um *UserModel) AddUser(newData Tbl_user) error {
	err := um.Connection.Create(&newData).Error
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) CekUser(hp string) bool {
	var data Tbl_user
	if err := um.Connection.Where("hp = ?", hp).First(&data).Error; err != nil {
		return false
	}
	return true
}

func (um *UserModel) Update(hp string, data Tbl_user) error {
	if err := um.Connection.Model(&data).Where("hp = ?", hp).Update("nama", data.Nama).Update("password", data.Password).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) GetAllUser() ([]Tbl_user, error) {
	var result []Tbl_user

	if err := um.Connection.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (um *UserModel) GetProfile(hp string) (Tbl_user, error) {
	var result Tbl_user
	if err := um.Connection.Where("hp = ?", hp).First(&result).Error; err != nil {
		return Tbl_user{}, err
	}
	return result, nil
}

func (um *UserModel) Login(hp string, password string) (Tbl_user, error) {
	var result Tbl_user
	if err := um.Connection.Where("hp = ? AND password = ?", hp, password).First(&result).Error; err != nil {
		return Tbl_user{}, err
	}
	return result, nil
}
