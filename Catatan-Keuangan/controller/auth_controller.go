package controller

import (
	"RestFullAPI/model/dto"
	"RestFullAPI/model/dto/common"
	"RestFullAPI/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC      usecase.AuthUsecase
	routerGroup *gin.RouterGroup
}

func (a *AuthController) loginHandler(ctx *gin.Context) {
	var payload dto.AuthReqDto
	err := ctx.ShouldBindJSON(&payload)
	log.Println("ini login", payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	token, err := a.authUC.Login(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "login server error")
	}
	common.SendSingleRespon(ctx, token, "Succes login")
}

func (a *AuthController) Route() {
	a.routerGroup.POST("/auth/login", a.loginHandler)
}

func NewAuthController(authUC usecase.AuthUsecase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUC: authUC, routerGroup: rg}
}
