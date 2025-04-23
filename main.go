package main

import (
	"github.com/gin-gonic/gin"
	"github/andro-kes/Blog/controllers"
	"github/andro-kes/Blog/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.DBMiddleWare)

	usersRouter := router.Group("/users")
	usersRouter.POST("/login", controllers.LoginHandler)
	usersRouter.POST("/signup", controllers.SignupHandler)
	usersRouter.POST("/reset_password", controllers.ResetPasswordHandler)
	usersRouter.POST("/logout", controllers.LogoutHandler)
}