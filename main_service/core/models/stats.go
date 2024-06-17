package models

type PostsTop struct {
	Type  string `uri:"type" binding:"required"`
	Count uint64 `uri:"count" binding:"required"`
}

type UsersTop struct {
	Count uint64 `uri:"count" binding:"required"`
}

type TopPostItem struct {
	PostId      uint64 `json:"post_id" binding:"required"`
	AuthorId    uint64 `json:"author_id" binding:"required"`
	AuthorLogin string `json:"author_login" binding:"required"`
	StatsNumber uint64 `json:"stats_number" binding:"required"`
}

type TopUserItem struct {
	UserId    uint64 `json:"user_id" binding:"required"`
	UserLogin string `json:"user_login" binding:"required"`
	SumLikes  uint64 `json:"sum_likes" binding:"required"`
}
