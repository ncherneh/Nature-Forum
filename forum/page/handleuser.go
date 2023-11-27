package page

import (
	"fmt"
	"forum/db"
	"forum/dop"
	"forum/structs"
	"html/template"
	"net/http"
	"strings"
)

type AccountToUser struct {
	AuthUser             *structs.User
	User                 *structs.User
	Posts                []*structs.Post
	Comments             []*structs.Comment
	Likes_Posts          []*structs.Post
	Likes_Comments       []*structs.Comment
	PreviewComments      []string
	PreviewLikesComments []string
}

// Handle requests on the user profile page
func HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	cookieUserID := dop.GetUserInCookie(r)
	username := strings.Replace(r.URL.Path, "/user/", "", 1)
	account := db.GetUserDataByUsername(username)
	if account == nil || username == "user" {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	tmpl, err := template.ParseFiles("./templates/user.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var reg AccountToUser
	if cookieUserID == nil || cookieUserID.Username != account.Username {
		account.Email = ""
	}
	reg.User = dop.GetUserInCookie(r)
	reg.AuthUser = account
	reg.Posts = db.GetPostsByUserId(account.ID)
	reg.Comments = db.GetCommentsByUserId(account.ID)
	reg.Likes_Posts = db.GetLikedPostsByUserID(account.ID)
	reg.Likes_Comments = db.GetLikedCommentsByUserID(account.ID)
	reg.Comments = db.GetCommentsByUserId(account.ID)
	reg.Likes_Comments = db.GetLikedCommentsByUserID(account.ID)

	// Generate preview for comments
	for _, comment := range reg.Comments {
		previewComment := generatePreviewContent(comment.Content, 100)
		reg.PreviewComments = append(reg.PreviewComments, previewComment)
	}

	// Generate preview for liked comments
	for _, likedComment := range reg.Likes_Comments {
		previewLikedComment := generatePreviewContent(likedComment.Content, 100)
		reg.PreviewLikesComments = append(reg.PreviewLikesComments, previewLikedComment)
	}

	err = tmpl.Execute(w, reg)
	if err != nil {
		fmt.Println(err)
	}
}

// Return an authorized user with a session cookie
func GetAuthUser(r *http.Request) *structs.User {
	cookie, err := r.Cookie("session")
	if err == nil {
		session := db.GetSessionData(cookie.Value)
		if session == nil {
			return nil
		}
		return db.GetUserDataById(session.UserID)
	}
	return nil
}
