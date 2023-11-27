package db

import (
	"database/sql"
	"fmt"
	"forum/structs"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB = OpenDB()

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("DB Operational")
	return db
}

// User
func InsertUserData(user structs.User) {
	_, err := Db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		user.Username, user.Email, user.Password, time.Now())
	if err != nil {
		fmt.Println(err)
	}
}

func GetUserDataByEmail(email string) *structs.User {
	var u structs.User
	err := Db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE email = ?", email).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get user data by email failed: " + err.Error())
		}
		return nil
	}
	return &u
}

func GetUserDataById(id int) *structs.User {
	var u structs.User
	err := Db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get user data by id failed: " + err.Error())
		}
		return nil
	}
	return &u
}

func GetUserDataByUsername(username string) *structs.User {
	var u structs.User
	err := Db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get user data by username failed: " + err.Error())
		}
		return nil
	}
	return &u
}

func UpdateUserData(userId int, newUsername, newEmail, newPassword string) {
	_, err := Db.Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?", newUsername, newEmail, newPassword, userId)
	if err != nil {
		fmt.Println(err)
	}
}

// Sessions-cookie
func InsertCookieData(session structs.Session) {
	_, err := Db.Exec("INSERT INTO sessions (id, user_id, expires) VALUES (?, ?, ?)", session.ID, session.UserID, session.Expires)
	if err != nil {
		fmt.Println("Inserting cookie failed: ", err.Error())
	}
}

func GetSessionData(sessionID string) *structs.Session {
	var s structs.Session
	err := Db.QueryRow("SELECT id, user_id, expires FROM sessions WHERE id = ?", sessionID).Scan(&s.ID, &s.UserID, &s.Expires)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get all post data failed: " + err.Error())
		}
		return nil
	}
	return &s
}

func RemoveCookieData(userID int) error {
	_, err := Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("failed to remove cookie data: " + err.Error())
	}
	return err
}

// Posts
func InsertPostData(post structs.Post) int {
	result, err := Db.Exec("INSERT INTO posts (user_id, title, content, image, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
	post.UserID, post.Title, post.Content, post.Image, post.Likes, post.Dislikes, post.CreatedAt)
    if err != nil {
        fmt.Println("Inserting post failed: " + err.Error())
        return 0
    }
    id, err := result.LastInsertId()
    if err != nil {
        fmt.Println("Getting last insert id failed: " + err.Error())
        return 0
    }
    for _, categoryID := range post.Category {
        _, err = Db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", id, categoryID)
        if err != nil {
            fmt.Println("Inserting post category mapping failed: " + err.Error())
            return 0
        }
    }
    return int(id)
}

func GetPost() []*structs.Post {
    var p []*structs.Post
	tabl, err := Db.Query("SELECT id, user_id, title, content, image, likes, dislikes, created_at FROM posts")
    if err != nil {
        fmt.Println("Getting all posts failed: " + err.Error())
    }
    defer tabl.Close()
    for tabl.Next() {
        var r structs.Post
        if err := tabl.Scan(&r.ID, &r.UserID, &r.Title, &r.Content, &r.Image, &r.Likes, &r.Dislikes, &r.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}		
        categoryRows, err := Db.Query("SELECT category_id FROM post_categories WHERE post_id = ?", r.ID)
        if err != nil {
            fmt.Println("Getting post categories failed: " + err.Error())
            continue
        }

        defer categoryRows.Close()

        var categories []int
        for categoryRows.Next() {
            var categoryID int
            if err := categoryRows.Scan(&categoryID); err != nil {
                fmt.Println("Error scanning category ID: " + err.Error())
                continue
            }
            categories = append(categories, categoryID)
        }
        r.Category = categories

        p = append(p, &r)
    }
    // fmt.Println("all posts", p)
    return p
}

func GetPostsByUserId(userId int) []*structs.Post {
	var p []*structs.Post
	tabl, err := Db.Query("SELECT id, user_id, title, content, image, likes, dislikes, created_at FROM posts WHERE user_id = ?", userId)
	if err != nil {
		fmt.Println("Getting posts by user id failed: " + err.Error())
		return nil
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.Post
		if err := tabl.Scan(&r.ID, &r.UserID, &r.Title, &r.Content, &r.Image, &r.Likes, &r.Dislikes, &r.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}		
		p = append(p, &r)
	}
	return p
}

func GetAllPostData(postId int) *structs.Post {
	var p structs.Post
	err := Db.QueryRow("SELECT p.id, p.user_id, p.title, p.content, p.image, p.likes, p.dislikes, p.created_at, u.username FROM posts AS p LEFT JOIN users AS u ON p.user_id = u.id WHERE p.id = ?", postId).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Image, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get all post data failed: " + err.Error())
		}
		return nil
	}
	return &p
}

func UpdatePostData(postId, likes, dislikes, comments int) {
	_, err := Db.Exec("UPDATE posts SET likes = ?, dislikes = ? WHERE id = ?;", likes, dislikes, postId)
	if err != nil {
		fmt.Println("Updating post failed: " + err.Error())
	}
}

func DeletePost(postID int) {
	_, err := Db.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		fmt.Println("Deleting post failed: " + err.Error())
	}
}

// Comments
func InsertComment(comment structs.Comment) {
	_, err := Db.Exec("INSERT INTO comments (post_id, user_id, content, likes, dislikes, created_at) VALUES(?, ?, ?, ?, ?, ?);",
		comment.PostID, comment.UserID, comment.Content, comment.Likes, comment.Dislikes, comment.CreatedAt)
	if err != nil {
		fmt.Println("Inserting comment failed: " + err.Error())
	}
}

func GetCommentsByPostId(postId int) []*structs.Comment {
	tabl, err := Db.Query("SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.likes, comments.dislikes, comments.created_at, users.username FROM comments LEFT JOIN users ON comments.user_id = users.id WHERE post_id = ?", postId)
	if err != nil {
		fmt.Println("Getting comments failed: " + err.Error())
		return nil
	}
	c := make([]*structs.Comment, 0)
	defer tabl.Close()
	for tabl.Next() {
		var r structs.Comment
		if err := tabl.Scan(&r.ID, &r.PostID, &r.UserID, &r.Content, &r.Likes, &r.Dislikes, &r.CreatedAt, &r.Username); err != nil {
			fmt.Println(err)
			continue
		}
		c = append(c, &r)
	}
	return c
}

func GetCommentsByUserId(userId int) []*structs.Comment {
	tabl, err := Db.Query("SELECT id, post_id, user_id, content, likes, dislikes, created_at FROM comments WHERE user_id = ?", userId)
	if err != nil {
		fmt.Println("Getting commentsByUserId failed: " + err.Error())
		return nil
	}
	c := make([]*structs.Comment, 0)
	defer tabl.Close()
	for tabl.Next() {
		var comment structs.Comment
		if err := tabl.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}
		c = append(c, &comment)
	}
	return c
}

func GetAllCommentData(commentId int) *structs.Comment {
	var c structs.Comment
	err := Db.QueryRow("SELECT id, post_id, user_id, content, likes, dislikes, created_at FROM comments WHERE id = ?", commentId).Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get all comment data failed: " + err.Error())
		}
		return nil
	}
	return &c
}

func UpdateCommentData(commentId int, newLikes int, newDislikes int) {
	_, err := Db.Exec("UPDATE comments SET likes = ?, dislikes = ? WHERE id = ?", newLikes, newDislikes, commentId)
	if err != nil {
		fmt.Println("Updating comment failed: " + err.Error())
	}
}

func RemoveComment(commentId int) {
	_, err := Db.Exec("DELETE FROM comments WHERE id = ?", commentId)
	if err != nil {
		fmt.Println("Removing comment failed: " + err.Error())
	}
}

// Likes & Dislikes Posts
func InsertPostLike(postLike structs.TotalLikesPost) {
	_, err := Db.Exec("INSERT OR IGNORE INTO total_likes_post (post_id, user_id) VALUES (?, ?)", postLike.PostID, postLike.UserID)
	if err != nil {
		fmt.Println("InsertPostLike failed: " + err.Error())
	}
}

func InsertPostDislike(postDislike structs.TotalDislikesPost) {
	_, err := Db.Exec("INSERT OR IGNORE INTO total_dislikes_post (post_id, user_id) VALUES (?, ?)", postDislike.PostID, postDislike.UserID)
	if err != nil {
		fmt.Println("InsertPostDislike failed: " + err.Error())
	}
}

func GetPostLikesByUserID(userID, postID int) structs.TotalLikesPost {
	var l structs.TotalLikesPost
	err := Db.QueryRow("SELECT post_id, user_id FROM total_likes_post WHERE user_id = ? AND post_id = ?", userID, postID).
		Scan(&l.PostID, &l.UserID)
	if err != nil {
		return structs.TotalLikesPost{}
	}
	return l
}

func GetPostDislikesByUserID(userID, postID int) structs.TotalDislikesPost {
	var d structs.TotalDislikesPost
	err := Db.QueryRow("SELECT post_id, user_id FROM total_dislikes_post WHERE user_id = ? AND post_id = ?", userID, postID).
		Scan(&d.PostID, &d.UserID)
	if err != nil {
		return structs.TotalDislikesPost{}
	}
	return d
}

func GetLikedPostsByUserID(userID int) []*structs.Post {
    rows, err := Db.Query("SELECT p.id, p.user_id, u.username, p.title, p.content, p.likes, p.dislikes, p.created_at FROM posts p JOIN total_likes_post lp ON p.id = lp.post_id JOIN users u ON p.user_id = u.id WHERE lp.user_id = ?", userID)
    if err != nil {
        fmt.Println("Error in GetLikedPostsByUserID: ", err)
        return nil
    }
    defer rows.Close()

    var posts []*structs.Post
    for rows.Next() {
        var p structs.Post
        err = rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.CreatedAt)
        if err != nil {
            fmt.Println("Error scanning row in GetLikedPostsByUserID: ", err)
            continue
        }
        posts = append(posts, &p)
    }
    return posts
}

func GetAllPostLikesByUserID(userID int) []*structs.TotalLikesPost {
	var l []*structs.TotalLikesPost
	tabl, err := Db.Query("SELECT post_id, user_id FROM total_likes_post WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllPostLikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.TotalLikesPost
		if err := tabl.Scan(&r.PostID, &r.UserID); err != nil {
			return nil
		}
		l = append(l, &r)
	}
	return l
}

func GetAllPostDislikesByUserID(userID int) []*structs.TotalDislikesPost {
	var d []*structs.TotalDislikesPost
	tabl, err := Db.Query("SELECT post_id, user_id FROM total_dislikes_post WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllPostDislikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.TotalDislikesPost
		if err := tabl.Scan(&r.PostID, &r.UserID); err != nil {
			return nil
		}
		d = append(d, &r)
	}
	return d
}

func RemovePostLikesByUserID(userID, postID int) error {
	_, err := Db.Exec("DELETE FROM total_likes_post WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func RemovePostDislikesByUserID(userID, postID int) error {
	_, err := Db.Exec("DELETE FROM total_dislikes_post WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// Likes & Dislikes Comments
func InsertCommentLike(commentLike structs.TotalLikesComment) {
	_, err := Db.Exec("INSERT OR IGNORE INTO total_likes_comment (comment_id, user_id) VALUES (?, ?)", commentLike.CommentID, commentLike.UserID)
	if err != nil {
		fmt.Println("InsertCommentLike failed: " + err.Error())
	}
}

func InsertCommentDislike(commentDislike structs.TotalDislikesComment) {
	_, err := Db.Exec("INSERT OR IGNORE INTO total_dislikes_comment (comment_id, user_id) VALUES (?, ?)", commentDislike.CommentID, commentDislike.UserID)
	if err != nil {
		fmt.Println("InsertCommentDislike failed: " + err.Error())
	}
}

func GetCommentLikesByUserID(userID, commentID int) structs.TotalLikesComment {
	var l structs.TotalLikesComment
	err := Db.QueryRow("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&l.CommentID, &l.UserID)
	if err != nil {
		return structs.TotalLikesComment{}
	}
	return l
}

func GetCommentDislikesByUserID(userID, commentID int) structs.TotalDislikesComment {
	var d structs.TotalDislikesComment
	err := Db.QueryRow("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&d.CommentID, &d.UserID)
	if err != nil {
		return structs.TotalDislikesComment{}
	}
	return d
}

func GetAllCommentLikesByUserID(userID int) []*structs.TotalLikesComment {
	var l []*structs.TotalLikesComment
	tabl, err := Db.Query("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllCommentLikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.TotalLikesComment
		if err := tabl.Scan(&r.CommentID, &r.UserID); err != nil {
			return nil
		}
		l = append(l, &r)
	}
	return l
}

func GetAllCommentDislikesByUserID(userID int) []*structs.TotalDislikesComment {
	var d []*structs.TotalDislikesComment
	tabl, err := Db.Query("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllCommentDislikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.TotalDislikesComment
		if err := tabl.Scan(&r.CommentID, &r.UserID); err != nil {
			return nil
		}
		d = append(d, &r)
	}
	return d
}

func GetLikedCommentsByUserID(userID int) []*structs.Comment {
	rows, err := Db.Query("SELECT c.id, c.user_id, c.post_id, c.content, c.likes, c.dislikes FROM comments c JOIN total_likes_comment lc ON c.id = lc.comment_id WHERE lc.user_id = ?", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []*structs.Comment
	for rows.Next() {
		var c structs.Comment
		err = rows.Scan(&c.ID, &c.UserID, &c.PostID, &c.Content, &c.Likes, &c.Dislikes)
		if err != nil {
			return nil
		}
		comments = append(comments, &c)
	}
	return comments
}

func RemoveCommentLikesByUserID(userID, commentID int) error {
	_, err := Db.Exec("DELETE FROM total_likes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveCommentDislikesByUserID(userID, commentID int) error {
	_, err := Db.Exec("DELETE FROM total_dislikes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// Categories
func GetCategories() []*structs.Category {
	var c []*structs.Category
	tabl, err := Db.Query("SELECT id, name FROM categories")
	if err != nil {
		fmt.Println("Getting all categories failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		var r structs.Category
		if err := tabl.Scan(&r.ID, &r.Name); err != nil {
			fmt.Println(err)
			continue
		}
		c = append(c, &r)
	}
	return c
}

func GetCategoriesByPostId(postId int) []*structs.Category {
    var categories []*structs.Category
    rows, err := Db.Query("SELECT categories.id, categories.name FROM post_categories INNER JOIN categories ON post_categories.category_id = categories.id WHERE post_categories.post_id = ?", postId)
    if err != nil {
        fmt.Println("Get categories by post id failed: " + err.Error())
        return nil
    }
    defer rows.Close()
    for rows.Next() {
        var r structs.Category
        if err := rows.Scan(&r.ID, &r.Name); err != nil {
            fmt.Println(err)
            continue
        }
        categories = append(categories, &r)
    }
    return categories
}