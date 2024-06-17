package controllers

import (
	"context"
	"log"

	"main_service_core/models"
	post_pb "main_service_core/post"
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

	post_grpc_resp, err := post_pb.Client.GetById(context.TODO(), &post_pb.GetByIdRequest{
		PostId: post_id.PostId,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
	}

	err = reaction.ReportView(post_id.PostId, post_grpc_resp.Post.AuthorId, user_id)
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

	post_grpc_resp, err := post_pb.Client.GetById(context.TODO(), &post_pb.GetByIdRequest{
		PostId: post_id.PostId,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
	}

	err = reaction.ReportLike(post_id.PostId, post_grpc_resp.Post.AuthorId, user_id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
	}

	c.JSON(200, gin.H{})
}
