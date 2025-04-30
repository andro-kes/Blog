package main

import (
	"log"

	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.DBMiddleWare())

	defer func() {
		sqlDB, err := middlewares.DB.DB()
		if err != nil {
			log.Fatalln("Ошибка при попытке закрыть базу данных: %м", err)
		}
		sqlDB.Close()
		log.Println("Соединение с базой данных разорвано")
	} ()

	usersRouter := router.Group("/users")
	usersRouter.POST("/login", controllers.LoginHandler)
	usersRouter.POST("/signup", controllers.SignupHandler)
	usersRouter.POST("/logout", controllers.LogoutHandler)
	usersRouter.POST("/refresh_token", controllers.RefreshTokenHandler)
	usersRouter.GET("/authYandex", controllers.AuthYandexRedirectHandler)
	usersRouter.GET("/loginYandexHandler", controllers.LoginYandexHandler)

	router.Run(":8000")
}