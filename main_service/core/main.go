package main

import (
	"fmt"
	"log"
	"main_service_core/controllers"
	"main_service_core/db_utils"
	"main_service_core/middlewares"
	pb "main_service_core/proto/post"
	"main_service_core/utils"

	gin "github.com/gin-gonic/gin"
)

func main() {
	err := db_utils.StartUpDB()
	for err != nil {
		log.Println(err, "retrying...")
		err = db_utils.StartUpDB()
	}

	err = pb.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.POST("/api/v1/user/sign-up", controllers.SignUp)
	r.POST("/api/v1/user/sign-in", controllers.SignIn)

	authed := r.Group("")
	authed.Use(middlewares.AuthMiddleware)
	authed.PUT("/api/v1/user/personal", controllers.UpdatePersonal)
	authed.POST("/api/v1/post", controllers.CreatePost)
	authed.PUT("/api/v1/post/:id", controllers.UpdatePost)
	authed.DELETE("/api/v1/post/:id", controllers.DeletePost)
	authed.GET("/api/v1/post/:id", controllers.GetPostById)
	authed.POST("/api/v1/post/pagination", controllers.GetPostPagination)

	listening_line := fmt.Sprintf(":%s", utils.GetenvSafe("MAIN_SERVICE_PORT"))
	log.Printf("listening at %s\n", listening_line)
	r.Run(listening_line)
}
