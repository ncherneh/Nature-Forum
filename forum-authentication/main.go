package main

import (
	"fmt"
	"forum/page"
	"log"
	"net/http"
)

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))
	style := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", style))
	http.HandleFunc("/register", page.HandleRegister)
	http.HandleFunc("/login", page.HandleLogin)
	http.HandleFunc("/logout", page.HandleLogOut)
	http.HandleFunc("/", page.HandleMain)
	http.HandleFunc("/user/", page.HandleUser)
	http.HandleFunc("/create-post", page.CreatePost)
	http.HandleFunc("/post/", page.HandlerPost)
	http.HandleFunc("/like-post", page.HandlerLikePost)
	http.HandleFunc("/dislike-post", page.HandlerDislikePost)
	http.HandleFunc("/delete-post", page.DeletePost)
	http.HandleFunc("/like-comment", page.HandlerLikeComment)
	http.HandleFunc("/dislike-comment", page.HandlerDislikeComment)
	http.HandleFunc("/comment", page.HandleComment)
	http.HandleFunc("/delete-comment", page.DeleteComment)
	http.HandleFunc("/edit-post", page.EditPost)
	http.HandleFunc("/edit-comment", page.EditComment)
	http.HandleFunc("/notifications", page.Notificaion)

	http.HandleFunc("/google/register/", page.RegisterGoogle)
	http.HandleFunc("/google/login/", page.LoginGoogle)
	http.HandleFunc("/logincallback", page.LoginCallBack)
	http.HandleFunc("/registercallback", page.RegisterCallBack)

	http.HandleFunc("/github/login/", page.GithubLogin)
	http.HandleFunc("/login/github/callback", page.GithubCallback)

	fmt.Println("Open http://localhost:8083\nPress Ctrl+C to exit")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Internal Server Error: %v", http.StatusInternalServerError)
	}
}
