package user

import (
	"Encode/helper"
	"Encode/middlewares"
	model "Encode/model/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Model model.UserModel
}

func (us *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirim kurang sesuai", nil))
		}

		var processInput model.Tbl_user
		processInput.Hp = input.Hp
		processInput.Nama = input.Nama
		processInput.Password = input.Password

		err = us.Model.AddUser(processInput) // ini adalah fungsi yang kita buat sendiri
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}

func (us *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			for _, val := range err.(validator.ValidationErrors) {
				fmt.Println(val.Error())
			}
		}

		result, err := us.Model.Login(input.Hp, input.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		token, err := middlewares.GenerateJWT(result.Hp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem, gagal memproses data", nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "selamat anda berhasil login", responseData))

	}
}

func (us *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("hp")
		var input model.Tbl_user
		err := c.Bind(&input)
		if err != nil {
			log.Println("masalah baca input:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		isFound := us.Model.CekUser(hp)

		if !isFound {
			return c.JSON(http.StatusNotFound,
				helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
		}

		err = us.Model.Update(hp, input)

		if err != nil {
			log.Println("masalah database :", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan saat update data", nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "data berhasil di update", nil))
	}
}

func (us *UserController) ListUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		listUser, err := us.Model.GetAllUser()
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", listUser))
	}
}

func (us *UserController) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hpFromToken = middlewares.DecodeToken(c.Get("user").(*jwt.Token))
		result, err := us.Model.GetProfile(hpFromToken)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}
