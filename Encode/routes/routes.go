package routes

import (
	"Encode/config"
	"Encode/controller/book"
	"Encode/controller/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.UserController, tc book.BookController) {
	userRoute(c, ctl)
	todoRoute(c, tc)
}

func userRoute(c *echo.Echo, ctl user.UserController) {
	c.POST("/users", ctl.Register()) // register -> umum (boleh diakses semua orang)
	c.POST("/login", ctl.Login())
	c.GET("/users", ctl.ListUser(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.GET("/profile", ctl.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/users/:hp", ctl.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}

func todoRoute(c *echo.Echo, tc book.BookController) {
	c.GET("/book", tc.GetBooksController())
	c.GET("book/:id", tc.GetBookController())
	c.POST("/book", tc.AddBookController(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/book/:todoID", tc.UpdateBookController(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/book/:id", tc.DeleteBookController(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
