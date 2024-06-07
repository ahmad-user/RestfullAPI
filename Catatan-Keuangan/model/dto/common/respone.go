package common

import (
	"RestFullAPI/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSingleRespon(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, &dto.SingleRespone{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func SendMyResponse(c *gin.Context, data []interface{}, paging dto.Paging, message string) {
	c.JSON(http.StatusOK, &dto.ManyResponse{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusOK, &dto.SingleRespone{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
	})
}
