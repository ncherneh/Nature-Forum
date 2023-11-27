package structs

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type TotalLikesPost struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
type TotalDislikesPost struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
type TotalLikesComment struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	CommentID int `json:"comment_id"`
}
type TotalDislikesComment struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	CommentID int `json:"comment_id"`
}
type Session struct {
	ID      int       `json:"id"`
	UserID  int       `json:"user_id"`
	Expires time.Time `json:"expires"`
}
type Post struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	Username             string    `json:"username"`
	Title                string    `json:"title"`
	Content              string    `json:"content"`
	Category             []int       `json:"category"`
	Comment              int       `json:"comment"`
	Likes                int       `json:"likes"`
	Dislikes             int       `json:"dislikes"`
	CreatedAt            time.Time `json:"created_at"`
	IsLikedByAuthUser    bool
	IsDislikedByAuthUser bool
}
type Comment struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	PostID               int       `json:"post_id"`
	Username             string    `json:"username"`
	Content              string    `json:"content"`
	Likes                int       `json:"likes"`
	Dislikes             int       `json:"dislikes"`
	CreatedAt            time.Time `json:"created_at"`
	IsLikedByAuthUser    bool
	IsDislikedByAuthUser bool
}