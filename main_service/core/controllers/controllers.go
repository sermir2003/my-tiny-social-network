package controllers

import (
	"log"
	db_utils "main_service_core/db_utils"
	jwt_utils "main_service_core/jwt_utils"
	"main_service_core/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if id == uuid.Nil {
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
	id, err := uuid.Parse(c.GetString("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	var data models.PersonalData
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Personal data is missing"})
		return
	}

	err = db_utils.UpdatePersonal(id, data)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{})
}
