package usecase

import (
	"RestFullAPI/model/dto"
	"RestFullAPI/usecase/service"
	"log"
)

type AuthUsecase interface {
	Login(payload dto.AuthReqDto) (dto.AuthResponDto, error)
}

type authUsecase struct {
	jwtService service.JwtService
	authorUC   UserUseCase
}

func (a *authUsecase) Login(payload dto.AuthReqDto) (dto.AuthResponDto, error) {
	author, err := a.authorUC.FindByUsername(payload.Username)
	log.Println("authoir", author)
	if err != nil {
		log.Println(err)
		return dto.AuthResponDto{}, err
	}
	token, err := a.jwtService.CreateToken(author)
	if err != nil {
		return dto.AuthResponDto{}, err
	}
	return token, err
}

func NewAuthUsecase(jwtService service.JwtService, authorUC UserUseCase) AuthUsecase {
	return &authUsecase{jwtService: jwtService, authorUC: authorUC}
}
