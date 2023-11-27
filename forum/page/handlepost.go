package page

import (
	"forum/db"
	"forum/dop"
	"forum/structs"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Handler for requests to pages of posts
func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/post/") {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	VisualPost(w, r)
}

// Visual the post page display
func VisualPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postname := strings.Replace(r.URL.Path, "/post/", "", 1)
	postIdInt, err := strconv.Atoi(postname)
	if err != nil {
		HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	postId := db.GetAllPostData(postIdInt)
	if postId == nil || postname == "post" {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	postCategories := db.GetCategoriesByPostId(postId.ID)
	p.Categories = postCategories

	authUser := GetAuthUser(r)
	if authUser == nil {
		authUser = &structs.User{}
	}
	p.AuntUser = authUser

	p.Post = postId
	postLikes := db.GetAllPostLikesByUserID(p.AuntUser.ID)
	if len(postLikes) > 0 {
		p.PostLikes = postLikes[0]
	}

	postDislikes := db.GetPostDislikesByUserID(p.AuntUser.ID, p.Post.ID)
	postDislikesList := []*structs.TotalDislikesPost{&postDislikes}
	if len(postDislikesList) > 0 && postDislikesList[0] != nil {
		p.PostDislikes = postDislikesList[0]
	}

	comments := db.GetCommentsByPostId(p.Post.ID)
	sortCommentsByDate(comments)

	for _, c := range comments {
		commentLike := db.GetCommentLikesByUserID(p.AuntUser.ID, c.ID)
		commentDislike := db.GetCommentDislikesByUserID(p.AuntUser.ID, c.ID)
		if commentLike.CommentID != 0 {
			c.IsLikedByAuthUser = true
		}
		if commentDislike.CommentID != 0 {
			c.IsDislikedByAuthUser = true
		}
	}
	users := dop.GetUserInCookie(r)
	isLoggedIn := users != nil

	var likedPosts []*structs.TotalLikesPost
	var dislikedPosts []*structs.TotalDislikesPost
	if isLoggedIn {
		likedPosts = db.GetAllPostLikesByUserID(users.ID)
		dislikedPosts = db.GetAllPostDislikesByUserID(users.ID)
	}

	isLiked := false
	for _, likePost := range likedPosts {
		if likePost.PostID == p.Post.ID {
			isLiked = true
		}
	}
	p.Post.IsLikedByAuthUser = isLiked

	isDisliked := false
	for _, dislikePost := range dislikedPosts {
		if dislikePost.PostID == p.Post.ID {
			isDisliked = true
		}
	}
	p.Post.IsDislikedByAuthUser = isDisliked

	type Posts_Comments struct {
		UserID     *structs.User
		IsLoggedIn bool
		*structs.Post
		Comments   []*structs.Comment
		AuthUserID int
		Categories []*structs.Category
	}

	p_c := &Posts_Comments{
		UserID:     users,
		IsLoggedIn: isLoggedIn,
		Post:       p.Post,
		Comments:   comments,
		AuthUserID: p.AuntUser.ID,
		Categories: postCategories,
	}
	tmpl, err := template.ParseFiles("./templates/post.html", "./templates/main.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	err = tmpl.ExecuteTemplate(w, "post.html", p_c)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// Handle the creation of a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/create_post.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	users := dop.GetUserInCookie(r)
	isLoggedIn := users != nil
	categories := Category()
	if r.Method == "GET" {
		currentUser := GetAuthUser(r)
		if currentUser == nil {
			HandleError(w, http.StatusNonAuthoritativeInfo, "Please, log in or register to have an opportunity to create a post")
			return
		}
		data := struct {
			Err        int
			Message    string
			Categories []*structs.Category
			UserID     *structs.User
			IsLoggedIn bool
		}{
			Err:        0,
			Message:    "",
			Categories: categories,
			UserID:     users,
			IsLoggedIn: isLoggedIn,
		}
		tmpl.Execute(w, data)
	} else if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			HandleError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		title := r.FormValue("title-create-post")
		content := r.FormValue("text-create-post")
		categoryIDsStr := r.Form["category-create-post[]"]
		if len(categoryIDsStr) == 0 {
			data := struct {
				Err        int
				Message    string
				Categories []*structs.Category
				UserID     *structs.User
				IsLoggedIn bool
			}{
				Err:        1,
				Message:    "Choose at least one category",
				Categories: categories,
				UserID:     users,
				IsLoggedIn: isLoggedIn,
			}
			tmpl.Execute(w, data)
			// HandleError(w, http.StatusBadRequest, "No categories were selected")
			return
		}
		var categoryIDs []int
		for _, categoryIDStr := range categoryIDsStr {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				data := struct {
					Err        int
					Message    string
					Categories []*structs.Category
					UserID     *structs.User
					IsLoggedIn bool
				}{
					Err:        1,
					Message:    "Choose at least one category",
					Categories: categories,
					UserID:     users,
					IsLoggedIn: isLoggedIn,
				}
				tmpl.Execute(w, data)
				// HandleError(w, http.StatusBadRequest, "Invalid category ID")
				return
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		currentUser := GetAuthUser(r)
		// if currentUser == nil {
		// 	data := struct {
		// 		Err        int
		// 		Message    string
		// 		Categories []*structs.Category
		// 		UserID     *structs.User
		// 		IsLoggedIn bool
		// 	}{
		// 		Err:        2,
		// 		Message:    "Please, log in or register to have an opportunity to create a post",
		// 		Categories: categories,
		// 		UserID:     users,
		// 		IsLoggedIn: isLoggedIn,
		// 	}
		// 	tmpl.Execute(w, data)
		// 	return
		// }
		newPost := structs.Post{
			UserID:    currentUser.ID,
			Title:     title,
			Content:   content,
			Category:  categoryIDs,
			Likes:     0,
			Dislikes:  0,
			CreatedAt: time.Now(),
		}
		if title == "" || content == "" {
			err := 4
			mes := "you can not create an empty post"
			data := struct {
				Err     int
				Message string
			}{
				Err:     err,
				Message: mes,
			}
			tmpl.Execute(w, data)
			return
		}
		tempID := db.InsertPostData(newPost)
		newPost.ID = tempID
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Visual the post creation page
func VisualCreatePost(w http.ResponseWriter, r *http.Request) {
	p.AuntUser = GetAuthUser(r)

	if p.AuntUser == nil {
		HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	tmpl, err := template.ParseFiles("./templates/post.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var categories []*structs.Category
	if p.Categories != nil {
		categories = append(categories, p.Categories...)
	} else {
		categories = db.GetCategories()
	}

	var posts []*structs.Post
	if p.Post != nil {
		posts = append(posts, p.Post)
	} else {
		posts = db.GetPostsByUserId(p.AuntUser.ID)
	}

	var comments []*structs.Comment
	if p.Comments != nil {
		comments = append(comments, p.Comments)
	} else {
		comments = db.GetCommentsByUserId(p.AuntUser.ID)
	}

	var postLikes []*structs.TotalLikesPost
	if p.PostLikes != nil {
		postLikes = append(postLikes, p.PostLikes)
	} else {
		postLikes = db.GetAllPostLikesByUserID(p.AuntUser.ID)
	}

	data := struct {
		AuntUser   *structs.User
		Categories []*structs.Category
		Posts      []*structs.Post
		Comments   []*structs.Comment
		PostLikes  []*structs.TotalLikesPost
	}{
		AuntUser:   p.AuntUser,
		Categories: categories,
		Posts:      posts,
		Comments:   comments,
		PostLikes:  postLikes,
	}

	err = tmpl.ExecuteTemplate(w, "post.html", data)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// Function for deleting posts
func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	postIDStr := r.FormValue("post-id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	currentUser := GetAuthUser(r)
	if currentUser == nil {
		return
	}

	post := db.GetAllPostData(postID)
	if post == nil || post.UserID != currentUser.ID {
		http.Error(w, "You are not allowed to delete this post", http.StatusForbidden)
		return
	}

	db.DeletePost(postID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Sort posts by creation date
func sortPostsByDate(posts []*structs.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
}

// Function to search for posts by title and content
func SearchPosts(posts []*structs.Post, searchQuery string, categoryFilters []int) []*structs.Post {
	var filteredPosts []*structs.Post

	for _, post := range posts {
		if len(categoryFilters) > 0 && !contains(post.Category, categoryFilters[0]) {
			continue
		}
		if strings.Contains(strings.ToLower(post.Title), strings.ToLower(searchQuery)) ||
			strings.Contains(strings.ToLower(post.Content), strings.ToLower(searchQuery)) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

// Function for generating a brief preview of the post content
func generatePreviewContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	return content[:maxChars] + "..."
}

// Sorts the array of posts by decreasing ID and returns the sorted array
func ShowFilterInPost(Posts []*structs.Post) []*structs.Post {
	sort.Slice(Posts, func(i, j int) bool {
		return Posts[i].ID > Posts[j].ID
	})
	return Posts
}
