package main

import (
	"log"
	controllers "main_service_core/controllers"
	db_utils "main_service_core/db_utils"
	jwt_utils "main_service_core/jwt_utils"
	middlewares "main_service_core/middlewares"

	gin "github.com/gin-gonic/gin"
)

func main() {
	err := db_utils.StartUpDB()
	if err != nil {
		log.Fatal(err)
	}
	jwt_utils.StartUpJWT()

	r := gin.Default()
	r.POST("/api/v1/user/sign-up", controllers.SignUp)
	r.POST("/api/v1/user/sign-in", controllers.SignIn)

	authed := r.Group("")
	authed.Use(middlewares.AuthMiddleware)
	authed.PUT("/api/v1/user/update-personal", controllers.UpdatePersonal)

	r.Run(":8081")
}
