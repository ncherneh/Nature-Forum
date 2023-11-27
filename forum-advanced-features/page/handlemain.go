package page

import (
	"errors"
	"fmt"
	"forum/db"
	"forum/dop"
	"forum/structs"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type resPost struct {
	AuntUser        *structs.User
	Categories      []*structs.Category
	Post            *structs.Post
	Comments        *structs.Comment
	PostLikes       *structs.TotalLikesPost
	PostDislikes    *structs.TotalDislikesPost
	CommentLikes    *structs.TotalLikesComment
	CommentDislikes *structs.TotalDislikesComment
}

var p resPost

type MainPageData struct {
	UserID     *structs.User
	Posts      []*PostWithComments
	Categories []*structs.Category
	IsLoggedIn bool
	NewNotification bool
}

// Handle requests to the main forum page
func HandleMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, errors.New("Page Not Found").Error())
		return
	}

	categoryIdsRaw := r.URL.Query()["categories"]
	searchQuery := r.URL.Query().Get("search")

	var categoryFilters []int
	for _, categoryIdRaw := range categoryIdsRaw {
		categoryFilter, err := strconv.Atoi(categoryIdRaw)
		if err != nil {
			HandleError(w, http.StatusBadRequest, err.Error())
			return
		}
		categoryFilters = append(categoryFilters, categoryFilter)
	}

	var posts []*structs.Post = db.GetPost()
	if len(categoryFilters) > 0 {
		posts = FilterByCategory(posts, categoryFilters)
	}
	if len(categoryFilters) > 0 || searchQuery != "" {
		posts = SearchPosts(posts, searchQuery, categoryFilters)
	}
	sortPostsByDate(posts)
	users := dop.GetUserInCookie(r)
	isLoggedIn := users != nil
	sortOption := r.URL.Query().Get("sort")
	if sortOption == "likes" {
		FilterByLikes(posts)
	} else if sortOption == "dislikes" {
		FilterByDislikes(posts)
	}

	var likedPosts []*structs.TotalLikesPost
	var dislikedPosts []*structs.TotalDislikesPost
	if isLoggedIn {
		likedPosts = db.GetAllPostLikesByUserID(users.ID)
		dislikedPosts = db.GetAllPostDislikesByUserID(users.ID)
	}
	postsWithComments := make([]*PostWithComments, len(posts))

	for i, post := range posts {
		previewContent := generatePreviewContent(post.Content, 100)
		categories := db.GetCategoriesByPostId(post.ID)
		comments := db.GetCommentsByPostId(post.ID)
		postLikes := len(db.GetAllPostLikesByUserID(post.UserID))
		postDislikes := len(db.GetAllPostDislikesByUserID(post.UserID))
		author := db.GetUserDataById(post.UserID)
		authUser := GetAuthUser(r)
		authUserID := 0
		if authUser != nil {
			authUserID = authUser.ID
		}
		isLiked := false
		for _, likePost := range likedPosts {
			if likePost.PostID == post.ID {
				isLiked = true
			}
		}
		post.IsLikedByAuthUser = isLiked
		var categoryNames []string
		for _, categoryId := range post.Category {
			category := GetCategoryById(categoryId)
			categoryNames = append(categoryNames, category.Name)
		}

		categoriesString := strings.Join(categoryNames, ", ")
		isDisliked := false
		for _, dislikePost := range dislikedPosts {
			if dislikePost.PostID == post.ID {
				isDisliked = true
			}
		}
		previewCategories := generatePreviewContent(categoriesString, 50)
		post.IsDislikedByAuthUser = isDisliked
		postsWithComments[i] = &PostWithComments{
			Post:            post,
			Author:          author,
			CreatedAt:       post.CreatedAt,
			Comments:        comments,
			PostLikes:       postLikes,
			PostDislikes:    postDislikes,
			PreviewContent:  previewContent,
			PreviewCategories: previewCategories,
			AuthUserID:      authUserID,
			Categories:      categories,
		}
	}
	resMap := make(template.FuncMap)
	resMap["Split"] = strings.Split

	data := MainPageData{
		UserID:     users,
		Posts:      postsWithComments,
		Categories: db.GetCategories(),
		IsLoggedIn: isLoggedIn,
	}
	if isLoggedIn {
		data.NewNotification = CountOfNewNotification(users.ID) > 0
	}
	tmpl, err := template.ParseFiles("./templates/main.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

// Return all categories from the database
func Category() []*structs.Category {
	categories := db.GetCategories()
	return categories
}

// Filter the array of posts by category and return the filtered array
func FilterByCategory(posts []*structs.Post, categoryIds []int) []*structs.Post {
	filteredPosts := make([]*structs.Post, 0)
	for _, post := range posts {
		for _, categoryId := range categoryIds {
			if contains(post.Category, categoryId) {
				filteredPosts = append(filteredPosts, post)
				break
			}
		}
	}
	return filteredPosts
}

// Filter of posts by likes
func FilterByLikes(posts []*structs.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Likes > posts[j].Likes
	})
}

// Filter of posts by dislikes
func FilterByDislikes(posts []*structs.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Dislikes > posts[j].Dislikes
	})
}

// Check whether a given slice contains a given element
func contains(slice []int, item int) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func GetCategoryById(id int) *structs.Category {
	categories := db.GetCategories()  
	for _, category := range categories {
		if category.ID == id {
			return category 
		}
	}
	return nil  
}
