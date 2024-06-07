package controller

import (
	"RestFullAPI/middleware"
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"RestFullAPI/model/dto/common"
	"RestFullAPI/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type expensesControler struct {
	expensesUseCase usecase.ExpensesUseCase
	router          *gin.RouterGroup
	autMiddleware   middleware.AuthMiddleware
}

func (a *expensesControler) getAllExpenses(ctx *gin.Context) {
	page, er := strconv.Atoi(ctx.Query("page"))
	size, er2 := strconv.Atoi(ctx.Query("size"))
	if er != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "conversi page gagal")
	}
	if er2 != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "conversi size gagal")
	}

	listData, paging, err := a.expensesUseCase.FindAllExpenses(page, size)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	var data []interface{}
	for _, b := range listData {
		data = append(data, b)
	}
	common.SendMyResponse(ctx, data, paging, "OK")
}
func (a *expensesControler) getByIdExpenses(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := a.expensesUseCase.FindByIdExpenses(id)
	if err != nil {
		panic(err)
	}
	log.Println("id", data)
	ctx.JSON(http.StatusOK, &dto.SingleRespone{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: "OK",
		},
		Data: data,
	})
}
func (a *expensesControler) getByTypeExpenses(ctx *gin.Context) {
	typ := ctx.Param("transaction_type")

	data, err := a.expensesUseCase.FindByTypeExpenses(typ)
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

func (a *expensesControler) updatedExpenses(ctx *gin.Context) {
	var up model.Expenses
	err := ctx.ShouldBindJSON(&up)
	log.Println("update", err)
	if err != nil {
		panic(err)
	}
	data, err := a.expensesUseCase.UpdatedExpenses(up)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, "error update")
		return
	}
	ctx.JSON(http.StatusOK, &dto.SingleRespone{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: "OK",
		},
		Data: data,
	})

}

func (a *expensesControler) insertExpenses(c *gin.Context) {
	var expense model.Expenses
	if err := c.ShouldBindJSON(&expense); err != nil {
		log.Println("insert", err)
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	createdExpense, err := a.expensesUseCase.InsertExpenses(expense)
	if err != nil {
		common.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendSingleRespon(c, createdExpense, "Ok")
}
func (a *expensesControler) RoutingExpenses() {
	a.router.GET("/expenses", a.autMiddleware.CheckToken("user"), a.getAllExpenses)
	a.router.GET("/expenses/:id", a.autMiddleware.CheckToken("user"), a.getByIdExpenses)
	a.router.GET("/expenses/type/:transaction_type", a.autMiddleware.CheckToken("user"), a.getByTypeExpenses)
	a.router.PUT("/expenses/:balance", a.autMiddleware.CheckToken("user"), a.updatedExpenses)
	a.router.POST("/expenses", a.insertExpenses)
}

func NewExpensesController(expensesUC usecase.ExpensesUseCase, rg *gin.RouterGroup, authMiddle middleware.AuthMiddleware) *expensesControler {
	return &expensesControler{
		expensesUseCase: expensesUC,
		router:          rg,
		autMiddleware:   authMiddle,
	}
}
