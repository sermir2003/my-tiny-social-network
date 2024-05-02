package controllers

import (
	"log"
	"main_service_core/db_utils"
	"main_service_core/jwt_utils"
	"main_service_core/models"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(400, gin.H{"error": "Credentials are missing"})
		return
	}

	id, err := db_utils.CreateNewUser(creds)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	if id == 0 {
		c.JSON(400, gin.H{"error": "Login has already been taken"})
		return
	}

	token, err := jwt_utils.CreateJWT(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func SignIn(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(400, gin.H{"error": "Credentials are missing"})
		return
	}

	is_valid, err := db_utils.CheckPassword(creds)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}
	if !is_valid {
		c.JSON(400, gin.H{"error": "Incorrect login or password"})
		return
	}

	id, err := db_utils.GetId(creds.Login)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	token, err := jwt_utils.CreateJWT(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func UpdatePersonal(c *gin.Context) {
	id := getUserId(c)
	if id == 0 {
		return
	}

	var personal models.PersonalData
	if err := c.BindJSON(&personal); err != nil {
		c.JSON(400, gin.H{"error": "Personal data is missing"})
		return
	}

	err := db_utils.UpdatePersonal(id, personal)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{})
}
