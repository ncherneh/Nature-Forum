package page

import (
	"fmt"
	"forum/db"
	"forum/structs"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type PostWithComments struct {
	Post              *structs.Post
	Author            *structs.User
	CreatedAt         time.Time
	Comments          []*structs.Comment
	PostLikes         int
	PostDislikes      int
	PreviewContent    string
	PreviewCategories string
	AuthUserID        int
	Categories        []*structs.Category
	
}

// Handle request to add a comment to the post
func HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	currentUser := GetAuthUser(r)
	if currentUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err := r.ParseForm()
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("comment_content")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post := db.GetAllPostData(postID)
	if post == nil {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if content == "" {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	comment := structs.Comment{
		UserID:    currentUser.ID,
		PostID:    postID,
		Content:   content,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now(),
		
	}
	db.InsertComment(comment)

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

//Function for editing comment
func EditComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	err := r.ParseForm()
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	currentUser := GetAuthUser(r)
	if currentUser == nil {
		return
	}

	comment_content := r.FormValue("comment-edit-content")

	comment := db.GetAllCommentData(commentID)
	if comment == nil || comment.UserID != currentUser.ID {
		http.Error(w, "You are not allowed to edit this comment" + strconv.Itoa(commentID),  http.StatusForbidden)

		return
	}

	db.EditCommentData(commentID, comment_content)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

// Handle request to remove a comment
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	currentUser := GetAuthUser(r)
	if currentUser == nil {
		HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err := r.ParseForm()
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	comment := db.GetAllCommentData(commentID)
	if comment == nil {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if !IsCommentOwner(currentUser, comment) {
		HandleError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	var postID int
	postID = comment.PostID

	db.RemoveComment(commentID)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

// Check if a user is the comment owner
func IsCommentOwner(user *structs.User, comment *structs.Comment) bool {
	return user.ID == comment.UserID
}

// Sort comments in descending order of creation date
func sortCommentsByDate(comments []*structs.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})
}
