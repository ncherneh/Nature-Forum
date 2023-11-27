package page

import (
	"forum/db"
	"forum/dop"
	"net/http"
)

var new AccountToUser

// Handle logout request
func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	user := dop.CheckCookieIntegrity(w, r)
	if user != nil {
		err := db.RemoveCookieData(user.ID)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}



