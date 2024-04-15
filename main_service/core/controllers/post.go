package controllers

import (
	"context"
	"log"
	"main_service_core/models"

	"github.com/gin-gonic/gin"

	pb "main_service_core/proto/post"
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
	req := pb.CreateRequest{
		AuthorId: author_id,
		Content:  post_content.Content,
	}

	resp, err := pb.Client.Create(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
	req := pb.UpdateRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
		Content:  post_content.Content,
	}

	resp, err := pb.Client.Update(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	switch resp.Access {
	case pb.AccessResult_SUCCESS:
		c.JSON(200, gin.H{})
	case pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case pb.AccessResult_NOT_FOUND:
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
	req := pb.DeleteRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
	}

	resp, err := pb.Client.Delete(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	switch resp.Access {
	case pb.AccessResult_SUCCESS:
		c.JSON(200, gin.H{})
	case pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case pb.AccessResult_NOT_FOUND:
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
	req := pb.GetByIdRequest{
		AuthorId: author_id,
		PostId:   post_id.PostId,
	}

	resp, err := pb.Client.GetById(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	switch resp.Access {
	case pb.AccessResult_SUCCESS:
		c.JSON(200, resp.Post)
	case pb.AccessResult_ACCESS_DENIED:
		c.JSON(403, gin.H{"error": "You do not have the authority to modify this post"})
	case pb.AccessResult_NOT_FOUND:
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
	req := pb.GetPaginationRequest{
		AuthorId: author_id,
		Offset:   *pagination.Offset,
		Limit:    *pagination.Limit,
	}

	resp, err := pb.Client.GetPagination(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"posts": resp.Posts})
}
