package middleware

import (
	"RestFullAPI/usecase/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type AuthHeader struct {
	AuthHeader string `header:"Authorization" required:"true"`
}

func (a *authMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var header AuthHeader
		ctx.ShouldBindHeader(&header)
		log.Println("========header", header.AuthHeader)
		token := strings.Replace(header.AuthHeader, "Bearer ", "", -1)
		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			log.Println(" validate gagal")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("author", claims["authorId"])
		validRole := false
		for _, role := range roles {
			if role == claims["role"] {
				validRole = true
				break
			}
		}
		if !validRole {
			log.Println(" invalid role")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
