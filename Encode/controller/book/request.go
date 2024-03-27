package book

type BookRequest struct {
	Judul   string `json:"judul" form:"judul" validate:"required,max=30,min=5"`
	Penulis string `json:"penulis" form:"penulis" validate:"required,max=30,min=5"`
	Genre   string `json:"genre" form:"genre" validate:"required,max=30,min=5"`
	Tahun   string `json:"tahun" form:"tahun" validate:"required,max=2000,min=2024"`
}
