package controllers

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserId(c *gin.Context) uint64 {
	id, err := strconv.ParseUint(c.GetString("id"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return 0
	}
	return id
}
