package controllers

import (
	"log"

	"main_service_core/models"
	"main_service_core/reaction"

	"github.com/gin-gonic/gin"
)

func PostView(c *gin.Context) {
	user_id := getUserId(c)
	if user_id == 0 {
		return
	}

	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	err := reaction.ReportView(post_id.PostId, user_id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
	}

	c.JSON(200, gin.H{})
}

func PostLike(c *gin.Context) {
	user_id := getUserId(c)
	if user_id == 0 {
		return
	}

	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	err := reaction.ReportLike(post_id.PostId, user_id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
	}

	c.JSON(200, gin.H{})
}
