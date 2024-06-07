package controller

import (
	"RestFullAPI/middleware"
	"RestFullAPI/model/dto"
	"RestFullAPI/model/dto/common"
	"RestFullAPI/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userControler struct {
	userUseCase   usecase.UserUseCase
	router        *gin.RouterGroup
	autMiddleware middleware.AuthMiddleware
}

func (a *userControler) listHandler(ctx *gin.Context) {
	page, er := strconv.Atoi(ctx.Query("page"))
	size, er2 := strconv.Atoi(ctx.Query("size"))
	if er != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "conversi page gagal")
	}
	if er2 != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "conversi size gagal")
	}

	listData, paging, err := a.userUseCase.FindAll(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	var data []interface{}
	for _, b := range listData {
		data = append(data, b)
	}
	common.SendMyResponse(ctx, data, paging, "OK")
}

func (a *userControler) getByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := a.userUseCase.FindById(id)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, &dto.SingleRespone{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: "OK",
		},
		Data: data,
	})
}

func (a *userControler) RoutingUser() {
	// tbl authors
	a.router.GET("/user", a.autMiddleware.CheckToken("user"), a.listHandler)
	a.router.GET("/user/:id", a.autMiddleware.CheckToken("user"), a.getByIdHandler)
}

func NewUserController(authorUc usecase.UserUseCase, rg *gin.RouterGroup, authMiddle middleware.AuthMiddleware) *userControler {
	return &userControler{
		userUseCase:   authorUc,
		router:        rg,
		autMiddleware: authMiddle,
	}
}
