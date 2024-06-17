package controllers

import (
	"context"
	"log"
	"main_service_core/db_utils"
	"main_service_core/models"
	stats_pb "main_service_core/stats"

	"github.com/gin-gonic/gin"
)

func PostStats(c *gin.Context) {
	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	ctx := context.TODO()
	req := stats_pb.GetPostStatsRequest{
		PostId: post_id.PostId,
	}

	resp, err := stats_pb.Client.GetPostStats(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{
		"post_id": post_id.PostId,
		"views":   resp.CntViews,
		"likes":   resp.CntLikes,
	})
}

func PostsTop(c *gin.Context) {
	var posts_top models.PostsTop
	if err := c.ShouldBindUri(&posts_top); err != nil {
		c.JSON(400, gin.H{"error": "The top type or the count is missing"})
		return
	}

	var reaction_type stats_pb.ReactionType
	switch posts_top.Type {
	case "view":
		reaction_type = stats_pb.ReactionType_VIEW
	case "like":
		reaction_type = stats_pb.ReactionType_LIKE
	}

	ctx := context.TODO()
	req := stats_pb.GetTopPostsRequest{
		Type:    reaction_type,
		TopSize: posts_top.Count,
	}

	resp, err := stats_pb.Client.GetTopPosts(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	top := []models.TopPostItem{}
	for _, grpc_item := range resp.Top {
		author_login, err := db_utils.GetLoginById(grpc_item.AuthorId)
		if err != nil {
			log.Println(err.Error())
			c.JSON(500, gin.H{"error": "Internal error"})
			return
		}
		final_item := models.TopPostItem{
			PostId:      grpc_item.PostId,
			AuthorId:    grpc_item.AuthorId,
			AuthorLogin: author_login,
			StatsNumber: grpc_item.StatsNumber,
		}
		top = append(top, final_item)
	}

	c.JSON(200, gin.H{"top": top})
}

func UsersTop(c *gin.Context) {
	var users_top models.UsersTop
	if err := c.ShouldBindUri(&users_top); err != nil {
		c.JSON(400, gin.H{"error": "The count is missing"})
		return
	}

	ctx := context.TODO()
	req := stats_pb.GetTopUsersRequest{
		TopSize: users_top.Count,
	}

	resp, err := stats_pb.Client.GetTopUsers(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	top := []models.TopUserItem{}
	for _, grpc_item := range resp.Top {
		user_login, err := db_utils.GetLoginById(grpc_item.UserId)
		if err != nil {
			log.Println(err.Error())
			c.JSON(500, gin.H{"error": "Internal error"})
			return
		}
		final_item := models.TopUserItem{
			UserId:    grpc_item.UserId,
			UserLogin: user_login,
			SumLikes:  grpc_item.SumLikes,
		}
		top = append(top, final_item)
	}

	c.JSON(200, gin.H{"top": top})
}
