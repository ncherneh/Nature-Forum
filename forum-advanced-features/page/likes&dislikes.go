package page

import (
	"forum/db"
	"forum/structs"
	"net/http"
	"strconv"
)

type like struct {
	ID     int    `json:"id"`
	MarkID string `json:"mark_id"`
	UserID int    `json:"user_id"`
	Mark   string `json:"mark"`
}

// Handle requests to like the post
func HandlerLikePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, _ := strconv.Atoi(postIDStr)
	authUser := GetAuthUser(r)
	if authUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	LikesPost(like{ID: postID, UserID: authUser.ID, MarkID: "posts", Mark: "like"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

// Handle requests to dislike the post
func HandlerDislikePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	postID, _ := strconv.Atoi(postIDStr)
	authUser := GetAuthUser(r)
	if authUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	LikesPost(like{ID: postID, UserID: authUser.ID, MarkID: "posts", Mark: "dislike"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

// Handle the data on likes and dislikes for the post
func LikesPost(r like) {
	likes := db.GetPostLikesByUserID(r.UserID, r.ID)
	dislikes := db.GetPostDislikesByUserID(r.UserID, r.ID)
	posts := db.GetAllPostData(r.ID)
	if r.Mark == "like" {
		if likes != (structs.TotalLikesPost{}) {
			db.UpdatePostData(r.ID, posts.Likes-1, posts.Dislikes, posts.Comment)
			db.RemovePostLikesByUserID(r.UserID, r.ID)
			return
		}
		if dislikes != (structs.TotalDislikesPost{}) {
			posts.Dislikes--
		}
		db.UpdatePostData(r.ID, posts.Likes+1, posts.Dislikes, posts.Comment)
		db.RemovePostDislikesByUserID(r.UserID, r.ID)
		db.InsertPostLike(structs.TotalLikesPost{UserID: r.UserID, PostID: r.ID})
	} else if r.Mark == "dislike" {
		if dislikes != (structs.TotalDislikesPost{}) {
			db.UpdatePostData(r.ID, posts.Likes, posts.Dislikes-1, posts.Comment)
			db.RemovePostDislikesByUserID(r.UserID, r.ID)
			return
		}
		if likes != (structs.TotalLikesPost{}) {
			posts.Likes--
		}
		db.UpdatePostData(r.ID, posts.Likes, posts.Dislikes+1, posts.Comment)
		db.RemovePostLikesByUserID(r.UserID, r.ID)
		db.InsertPostDislike(structs.TotalDislikesPost{UserID: r.UserID, PostID: r.ID})
	}
}

// Handle requests to like the comment
func HandlerLikeComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Query().Get("comment_id")
	commentID, _ := strconv.Atoi(commentIDStr)
	authUser := GetAuthUser(r)
	if authUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	LikesComment(like{ID: commentID, UserID: authUser.ID, MarkID: "comments", Mark: "like"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

// Handle requests to dislike the comment
func HandlerDislikeComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Query().Get("comment_id")
	commentID, _ := strconv.Atoi(commentIDStr)
	authUser := GetAuthUser(r)
	if authUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	LikesComment(like{ID: commentID, UserID: authUser.ID, MarkID: "comments", Mark: "dislike"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

// Handle the data on likes and dislikes for the comment
func LikesComment(r like) {
	likes := db.GetCommentLikesByUserID(r.UserID, r.ID)
	dislikes := db.GetCommentDislikesByUserID(r.UserID, r.ID)
	comments := *db.GetAllCommentData(r.ID)
	if r.Mark == "like" {
		if likes != (structs.TotalLikesComment{}) {
			db.UpdateCommentData(r.ID, comments.Likes-1, comments.Dislikes)
			db.RemoveCommentLikesByUserID(r.UserID, r.ID)
			return
		}
		if dislikes != (structs.TotalDislikesComment{}) {
			comments.Dislikes--
		}
		db.UpdateCommentData(r.ID, comments.Likes+1, comments.Dislikes)
		db.RemoveCommentDislikesByUserID(r.UserID, r.ID)
		db.InsertCommentLike(structs.TotalLikesComment{UserID: r.UserID, CommentID: r.ID})

	} else if r.Mark == "dislike" {
		if dislikes != (structs.TotalDislikesComment{}) {
			db.UpdateCommentData(r.ID, comments.Likes, comments.Dislikes-1)
			db.RemoveCommentDislikesByUserID(r.UserID, r.ID)
			return
		}
		if likes != (structs.TotalLikesComment{}) {
			comments.Likes--
		}
		db.UpdateCommentData(r.ID, comments.Likes, comments.Dislikes+1)
		db.RemoveCommentLikesByUserID(r.UserID, r.ID)
		db.InsertCommentDislike(structs.TotalDislikesComment{UserID: r.UserID, CommentID: r.ID})
	}
}
