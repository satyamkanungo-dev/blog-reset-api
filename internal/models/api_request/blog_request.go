package apirequest

type BlogRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// delete multiple blogs
type DeleteBlogRequest struct {
	Ids []string `json:"id"`
}
