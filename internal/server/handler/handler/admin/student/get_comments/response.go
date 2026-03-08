package get_comments

type Comment struct {
	ID             int64  `json:"id"`
	Text           string `json:"text"`
	AuthorFullName string `json:"author_full_name"`
	CreatedAt      string `json:"created_at"`
}

type Response struct {
	Comments      []Comment `json:"comments"`
	CommentsCount int64     `json:"comments_count"`
}
