package models

import "time"

type PostContent struct {
	Content string `json:"content" binding:"required"`
}

// TODO: is needed?
type PostFull struct {
	PostId          uint32    `json:"post_id" binding:"required"`
	AuthorId        uint32    `json:"author_id" binding:"required"`
	Content         string    `json:"content" binding:"required"`
	CreateTimestamp time.Time `json:"create_timestamp" binding:"required"`
	UpdateTimestamp time.Time `json:"update_timestamp" binding:"required"`
}

type Pagination struct {
	Offset *uint32 `json:"offset" binding:"required"`
	Limit  *uint32 `json:"limit" binding:"required"`
}

type PostId struct {
	PostId uint32 `uri:"id" binding:"required"`
}
