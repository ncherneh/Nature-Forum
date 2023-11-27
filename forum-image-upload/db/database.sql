BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER,
	"username"	TEXT NOT NULL,
	"email"	TEXT NOT NULL,
	"password"	TEXT NOT NULL,
	"created_at"	DATETIME NOT NULL DEFAULT 'NOW',
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "categories" (
	"id"	INTEGER,
	"name"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "total_likes_post" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"post_id"	INTEGER,
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "posts"("id")
);
CREATE TABLE IF NOT EXISTS "total_dislikes_post" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"post_id"	INTEGER,
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "total_likes_comment" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"comment_id"	INTEGER,
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "total_dislikes_comment" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"comment_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "sessions" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"expires"	DATETIME NOT NULL DEFAULT 'NOW',
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"post_id"	INTEGER,
	"content"	TEXT NOT NULL,
	"likes"	INTEGER DEFAULT 0,
	"dislikes"	INTEGER DEFAULT 0,
	"created_at"	DATETIME NOT NULL DEFAULT 'NOW',
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("dislikes") REFERENCES "total_dislikes_comment"("id"),
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("likes") REFERENCES "total_likes_comment"("id")
);
CREATE TABLE IF NOT EXISTS "post_categories" (
	"post_id"	INTEGER,
	"category_id"	INTEGER,
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("category_id") REFERENCES "categories"("id")
);
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER,
	"user_id"	INTEGER,
	"title"	TEXT NOT NULL,
	"content"	TEXT NOT NULL,
	"image"	TEXT,
	"likes"	INTEGER DEFAULT 0,
	"dislikes"	INTEGER DEFAULT 0,
	"comment"	INTEGER,
	"created_at"	DATETIME NOT NULL DEFAULT 'NOW',
	FOREIGN KEY("likes") REFERENCES "total_likes_post"("id"),
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment") REFERENCES "comments"("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id"),
	FOREIGN KEY("dislikes") REFERENCES "total_dislikes_post"("id")
);