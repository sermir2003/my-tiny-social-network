package controllers

import (
	"context"
	"log"
	"main_service_core/models"

	"github.com/gin-gonic/gin"

	post_pb "main_service_core/proto/post"
)

func CreatePost(c *gin.Context) {
	author_id := getUserId(c)
	if author_id == 0 {
		return
	}

	var post_content models.PostContent
	if err := c.BindJSON(&post_content); err != nil {
		c.JSON(400, gin.H{"error": "Post content is missing"})
		return
	}

	ctx := context.TODO()
	req := post_pb.CreateRequest{
		AuthorId: author_id,
		Content:  post_content.Content,
	}

	resp, err := post_pb.Client.Create(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}
	c.JSON(200, gin.H{"post_id": resp.PostId})
}

func UpdatePost(c *gin.Context) {
	author_id := getUserId(c)
	if author_id == 0 {
		return
	}

	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	var post_content models.PostContent
	if err := c.BindJSON(&post_content); err != nil {
		c.JSON(400, gin.H{"error": "Post content is missing"})
		return
	}

	ctx := context.TODO()
	req := post_pb.UpdateRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
		Content:  post_content.Content,
	}

	resp, err := post_pb.Client.Update(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	switch resp.Access {
	case post_pb.AccessResult_SUCCESS:
		c.JSON(200, gin.H{})
	case post_pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case post_pb.AccessResult_NOT_FOUND:
		c.JSON(404, gin.H{"error": "There is no post with provided id"})
	default:
		log.Fatal("Unknown access result")
	}
}

func DeletePost(c *gin.Context) {
	author_id := getUserId(c)
	if author_id == 0 {
		return
	}

	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	ctx := context.TODO()
	req := post_pb.DeleteRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
	}

	resp, err := post_pb.Client.Delete(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	switch resp.Access {
	case post_pb.AccessResult_SUCCESS:
		c.JSON(200, gin.H{})
	case post_pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case post_pb.AccessResult_NOT_FOUND:
		c.JSON(404, gin.H{"error": "There is no post with provided id"})
	default:
		log.Fatal("Unknown access result")
	}
}

func GetPostById(c *gin.Context) {
	author_id := getUserId(c)
	if author_id == 0 {
		return
	}

	var post_id models.PostId
	if err := c.ShouldBindUri(&post_id); err != nil {
		c.JSON(400, gin.H{"error": "Post id is missing"})
		return
	}

	ctx := context.TODO()
	req := post_pb.GetByIdRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
	}

	resp, err := post_pb.Client.GetById(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	switch resp.Access {
	case post_pb.AccessResult_SUCCESS:
		c.JSON(200, resp.Post)
	case post_pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case post_pb.AccessResult_NOT_FOUND:
		c.JSON(404, gin.H{"error": "There is no post with provided id"})
	default:
		log.Fatal("Unknown access result")
	}
}

func GetPostPagination(c *gin.Context) {
	author_id := getUserId(c)
	if author_id == 0 {
		return
	}

	var pagination models.Pagination
	if err := c.BindJSON(&pagination); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "Pagination parameters are missing"})
		return
	}

	ctx := context.TODO()
	req := post_pb.GetPaginationRequest{
		AuthorId: author_id,
		Offset:   *pagination.Offset,
		Limit:    *pagination.Limit,
	}

	resp, err := post_pb.Client.GetPagination(ctx, &req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(200, gin.H{"posts": resp.Posts})
}
