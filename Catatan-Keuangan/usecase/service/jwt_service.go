package service

import (
	"RestFullAPI/config"
	"RestFullAPI/model"
	"RestFullAPI/model/dto"
	"RestFullAPI/utils"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(author model.Users) (dto.AuthResponDto, error)
	ValidateToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	co config.TokenConfig
}

// CreateToken implements JwtService.
func (j *jwtService) CreateToken(author model.Users) (dto.AuthResponDto, error) {
	claim := utils.CustomeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.co.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.co.ExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:     author.Role,
		AuthorId: author.Id,
	}
	token := jwt.NewWithClaims(j.co.SigningMethod, claim)
	ss, err := token.SignedString(j.co.SignatureKey)
	if err != nil {
		return dto.AuthResponDto{}, fmt.Errorf("failed created access token")
	}
	return dto.AuthResponDto{Token: ss}, nil
}

// ValidateToken implements JwtService.
func (j *jwtService) ValidateToken(tokenHeader string) (jwt.MapClaims, error) {
	log.Println("kepanggil")
	// log.Println("ValidateToken dipanggil dengan token:", tokenHeader)
	// parse, err := jwt.Parse(tokenHeader, func(parse *jwt.Token) (interface{}, error) {
	// 	log.Println(parse)
	// 	return j.co.SignatureKey, nil
	// })
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return j.co.SignatureKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to verify token")
	}
	return claims, nil
}

func NewJwtService(c config.TokenConfig) JwtService {
	return &jwtService{co: c}
}
