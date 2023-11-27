package page

import (
	"forum/db"
	"forum/structs"
	"net/http"
	"sort"
	"text/template"
	"time"
)

type MainData struct {
	Activity   []*structs.PostActivity
	IsLoggedIn bool
	User       *structs.User
}

func Notificaion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	currentUser := GetAuthUser(r)
	if currentUser == nil {
		HandleError(w, http.StatusNonAuthoritativeInfo, "Please, log in or register to have an opportunity to see your notifications")
		return
	}
	tmpl, err := template.ParseFiles("./templates/notification.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	var post_activ []*structs.PostActivity
	
	post_activ = append(post_activ, db.GetPostsActivityLikes(currentUser.ID)...)
	post_activ = append(post_activ, db.GetPostsActivityDislikes(currentUser.ID)...)
	post_activ = append(post_activ, db.GetPostsActivityComments(currentUser.ID)...)
	sort.Slice(post_activ, func(i, j int) bool {
		return post_activ[i].DataTime.After(post_activ[j].DataTime)
	})
	nofic := db.GetNotification(currentUser.ID)

	if nofic.UserID != currentUser.ID {
		db.InsertNotifications(currentUser.ID, time.Now())
	} else {
		db.UpdateNotification(currentUser.ID, time.Now())
	}
	Data := MainData{
		Activity:   post_activ,
		User:       currentUser,
		IsLoggedIn: true,
	}
	err = tmpl.ExecuteTemplate(w, "notification.html", Data)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

}
func CountOfNewNotification(user_id int) int {
	last_visited := db.GetNotification(user_id).Lastvisit

	var post_activ []*structs.PostActivity

	post_activ = append(post_activ, db.GetPostsActivityLikes(user_id)...)
	post_activ = append(post_activ, db.GetPostsActivityDislikes(user_id)...)
	post_activ = append(post_activ, db.GetPostsActivityComments(user_id)...)

	result := 0

	for _, p_a := range post_activ {
		if p_a.DataTime.After(last_visited) {
			result++
		}
	}
	return result
}
