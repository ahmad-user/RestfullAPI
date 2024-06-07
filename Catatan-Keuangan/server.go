package main

import (
	"database/sql"
	"fmt"

	"RestFullAPI/config"
	"RestFullAPI/controller"
	"RestFullAPI/middleware"
	"RestFullAPI/repository"
	"RestFullAPI/usecase"
	"RestFullAPI/usecase/service"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Server struct {
	userUC     usecase.UserUseCase
	authUC     usecase.AuthUsecase
	jwtService service.JwtService
	expensesUC usecase.ExpensesUseCase
	engine     *gin.Engine
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	// token
	authMiddle := middleware.NewAuthMiddleware(s.jwtService)
	// user
	controller.NewUserController(s.userUC, rg, authMiddle).RoutingUser()
	// expenses
	controller.NewExpensesController(s.expensesUC, rg, authMiddle).RoutingExpenses()
	//auth
	controller.NewAuthController(s.authUC, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	s.engine.Run(":2000")
}

func NewServer() *Server {
	c, _ := config.NewConfig()
	urlConnect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.DbPort, c.DbUser, c.DbPassword, c.DbName)

	database, err := sql.Open(c.Driver, urlConnect)
	if err != nil {
		panic("error connected database")
	}
	fmt.Println("connected database")

	authorRepo := repository.NewUserRepo(database)
	userUC := usecase.UserUseCase(authorRepo)

	expensesRepo := repository.NewExpensesRepo(database)
	expensesUC := usecase.ExpensesUseCase(expensesRepo)

	jwtService := service.NewJwtService(c.TokenConfig)
	authUseCase := usecase.NewAuthUsecase(jwtService, userUC)

	return &Server{
		userUC:     userUC,
		expensesUC: expensesUC,
		authUC:     authUseCase,
		jwtService: jwtService,
		engine:     gin.Default(),
	}
}
