package models

import "time"

type PostContent struct {
	Content string `json:"content" binding:"required"`
}

type PostFull struct {
	PostId          uint64    `json:"post_id" binding:"required"`
	AuthorId        uint64    `json:"author_id" binding:"required"`
	Content         string    `json:"content" binding:"required"`
	CreateTimestamp time.Time `json:"create_timestamp" binding:"required"`
	UpdateTimestamp time.Time `json:"update_timestamp" binding:"required"`
}

type Pagination struct {
	Offset *uint64 `json:"offset" binding:"required"`
	Limit  *uint32 `json:"limit" binding:"required"`
}

type PostId struct {
	PostId uint64 `uri:"id" binding:"required"`
}
