package page

import (
	"forum/db"
	"forum/dop"
	"net/http"
	"text/template"
)

// Handle requests to display and send a login form
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if r.Method == "GET" {
		if r.URL.Path != "/login" {
			HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		tmpl.Execute(w, "")
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		userFromDb := db.GetUserDataByEmail(email)
		if userFromDb != nil && dop.CheckPasswordHash(password, userFromDb.Password) {
			cookie := dop.GetCookies(*userFromDb)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			if err != nil {
				panic(err)
			}
			err := "Invalid email or password"
			tmpl.Execute(w, err)
			return
		}
	} else {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}