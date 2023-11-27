package main

import (
	"fmt"
	"forum/page"
	"log"
	"net/http"
	"time"
	"sync"
	"encoding/json"
	"net"
	"golang.org/x/time/rate"
)

type Message struct {
    Status string `json:"status"`
    Body   string `json:"body"`
}

func main() {
	mux := http.NewServeMux()
	style := http.FileServer(http.Dir("./static"))
	
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))
	mux.Handle("/static/", http.StripPrefix("/static", style))
	mux.Handle("/register", perClientRateLimiter(page.HandleRegister))
	mux.Handle("/login", perClientRateLimiter(page.HandleLogin))
	mux.Handle("/logout", perClientRateLimiter(page.HandleLogOut))
	mux.Handle("/", perClientRateLimiter(page.HandleMain))
	mux.Handle("/user/", perClientRateLimiter(page.HandleUser))
	mux.Handle("/create-post", perClientRateLimiter(page.CreatePost))
	mux.Handle("/post/", perClientRateLimiter(page.HandlerPost))
	mux.Handle("/like-post", perClientRateLimiter(page.HandlerLikePost))
	mux.Handle("/dislike-post", perClientRateLimiter(page.HandlerDislikePost))
	mux.Handle("/delete-post", perClientRateLimiter(page.DeletePost))
	mux.Handle("/like-comment", perClientRateLimiter(page.HandlerLikeComment))
	mux.Handle("/dislike-comment", perClientRateLimiter(page.HandlerDislikeComment))
	mux.Handle("/comment", perClientRateLimiter(page.HandleComment))
	mux.Handle("/delete-comment", perClientRateLimiter(page.DeleteComment))

	fmt.Println("Open https://localhost:8083\nPress Ctrl+C to exit")

	certificate := "security/domain.crt"
	key := "security/domain.key"

	server := &http.Server{
		Addr:         ":8083",
		Handler:      mux,
		IdleTimeout:  200 * time.Second,
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	}

	err := server.ListenAndServeTLS(certificate, key)
    if err != nil {
        log.Fatal("Internal Server Error: ", err)
    }
}

func perClientRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
    type client struct {
        limiter  *rate.Limiter
        lastSeen time.Time
    }
    var (
        mu      sync.Mutex
        clients = make(map[string]*client)
    )
    go func() {
        for {
            time.Sleep(time.Minute)
            mu.Lock()
            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(clients, ip)
                }
            }
            mu.Unlock()
        }
    }()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        mu.Lock()
        if _, found := clients[ip]; !found {
            clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
        }
        clients[ip].lastSeen = time.Now()
        if !clients[ip].limiter.Allow() {
            mu.Unlock()
            message := Message{
                Status: "Request Failed",
                Body:   "The API is at capacity, try again later.",
            }
            w.WriteHeader(http.StatusTooManyRequests)
            json.NewEncoder(w).Encode(&message)
            return
        }
        mu.Unlock()
        next(w, r)
    })
}