package book

import (
	"Encode/helper"
	"Encode/middlewares"
	"Encode/model/book"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type BookController struct {
	Model book.BookModel
}

func (g *BookController) GetBooksController() echo.HandlerFunc {
	return func(c echo.Context) error {
		book, err := g.Model.GetAllBooks()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Gagal mengambil pengguna",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Berhasil mengambil semua data books",
			"users":   book,
		})
	}
}

func (gl *BookController) GetBookController() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "ID pengguna tidak valid",
			})
		}
		user, err := gl.Model.GetAllById(userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Gagal mengambil data pengguna",
			})
		}
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"code":    http.StatusNotFound,
				"message": "Pengguna tidak ditemukan",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"Message": "Berhasil menambahkan book",
			"user":    user,
		})
	}
}

func (tc *BookController) AddBookController() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input BookRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		// Cek middleware (extract token)
		// c.Get("user").(*jwt.Token) -> notasi PASTI kalo mau mengambil jwt token pada echo

		hp := middlewares.DecodeToken(c.Get("user").(*jwt.Token))

		if hp == "" {
			log.Println("error decode token:", "hp tidak ditemukan")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "tidak dapat mengakses fitur ini", nil))
		}

		var inputProcess book.Tbl_book
		inputProcess.Judul = input.Judul
		inputProcess.Penulis = input.Penulis
		inputProcess.Genre = input.Genre
		inputProcess.Tahun = input.Tahun
		inputProcess.Pemilik = hp

		result, err := tc.Model.Insert(inputProcess)

		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada proses server", nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan kegiatan", result))
	}
}

func (tc *BookController) UpdateBookController() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input BookRequest

		readID := c.Param("todoID")
		fmt.Println(readID)
		cnv, err := strconv.Atoi(readID)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}
		err = c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		hp := middlewares.DecodeToken(c.Get("user").(*jwt.Token))

		if hp == "" {
			log.Println("error decode token:", "hp tidak ditemukan")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "tidak dapat mengakses fitur ini", nil))
		}

		var inputProcess book.Tbl_book
		inputProcess.Judul = input.Judul
		inputProcess.Penulis = input.Penulis
		inputProcess.Genre = input.Genre
		inputProcess.Tahun = input.Tahun
		inputProcess.Pemilik = hp

		result, err := tc.Model.UpdateKegiatan(hp, uint(cnv), inputProcess)

		if err != nil {
			log.Println("error update db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada proses server", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "berhasil mengubah data kegiatan", result))
	}

}

func (del *BookController) DeleteBookController() echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.Param("email")
		bookToDelete := book.Tbl_book{Pemilik: email}
		err := del.Model.DeleteBook(bookToDelete)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Gagal menghapus pengguna",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Berhasil menghapus pengguna",
		})
	}
}
