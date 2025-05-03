package main

import (
	"log"

	posts_controllers "github.com/andro-kes/Blog/controllers/posts"
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
	usersRouter.POST("/login", users_controllers.LoginHandler)
	usersRouter.POST("/signup", users_controllers.SignupHandler)
	usersRouter.POST("/logout", users_controllers.LogoutHandler)
	usersRouter.POST("/refresh_token", users_controllers.RefreshTokenHandler)
	usersRouter.GET("/authYandex", users_controllers.AuthYandexRedirectHandler)
	usersRouter.GET("/loginYandexHandler", users_controllers.LoginYandexHandler)

	postsRouter := router.Group("/posts")
	postsRouter.POST("/create", posts_controllers.CreatePostHandler)
	postsRouter.GET("/:id", posts_controllers.RetrievePostHandler)

	router.Run(":8000")
}