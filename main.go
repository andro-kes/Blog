package main

import (
	"log"

	"github.com/andro-kes/Blog/controllers"
	"github.com/andro-kes/Blog/middlewares"
	"github.com/andro-kes/Blog/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	
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

	router.Run(":8000")
}