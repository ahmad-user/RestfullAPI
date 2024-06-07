package book

import (
	"errors"

	"gorm.io/gorm"
)

type Tbl_book struct {
	gorm.Model
	Judul   string `json:"judul" form:"judul" validate:"required,max=30,min=5"`
	Penulis string `json:"penulis" form:"penulis" validate:"required,max=30,min=5"`
	Genre   string `json:"genre" form:"genre" validate:"required,max=30,min=5"`
	Tahun   string `json:"tahun" form:"tahun" validate:"required,max=4,min=2"`
	Pemilik string `gorm:"type:varchar(13);"`
}

type BookModel struct {
	Connection *gorm.DB
}

func (g *BookModel) GetAllBooks() ([]Tbl_book, error) {
	var book []Tbl_book
	if err := g.Connection.Model(&Tbl_book{}).Find(&book).Error; err != nil {
		return nil, err
	}
	return book, nil
}
func (gl *BookModel) GetAllById(bookID int) ([]Tbl_book, error) {
	var bookById []Tbl_book
	if err := gl.Connection.Model(&Tbl_book{}).First(&bookById, bookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return bookById, nil
}

func (tm *BookModel) Insert(kegiatanBaru Tbl_book) (Tbl_book, error) {
	if err := tm.Connection.Create(&kegiatanBaru).Error; err != nil {
		return Tbl_book{}, err
	}

	return kegiatanBaru, nil
}

func (tm *BookModel) ListKegiatan(pemilik string) ([]Tbl_book, error) {
	var result []Tbl_book
	if err := tm.Connection.Where("pemilik = ?", pemilik).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (tm *BookModel) UpdateKegiatan(pemilik string, todoID uint, data Tbl_book) (Tbl_book, error) {
	var qry = tm.Connection.Where("pemilik = ? AND id = ?", pemilik, todoID).Updates(data)
	if err := qry.Error; err != nil {
		return Tbl_book{}, err
	}

	if qry.RowsAffected < 1 {
		return Tbl_book{}, errors.New("no data affected")
	}

	return data, nil
}

func (del *BookModel) DeleteBook(delID Tbl_book) error {
	var data = del.Connection.Delete(&delID)
	if err := data.Error; err != nil {
		return err
	}
	if data.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
