package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rama-kairi/blog-api-golang-gin/controllers"
)

func main() {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Close the database connection

	// Create the table if it doesn't exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS blog (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		body TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}

	r := gin.Default()
	newBlogC := controllers.NewBlogStore(db)

	r.GET("/blog", newBlogC.GetAllBlogs)
	r.GET("/blog/:id", newBlogC.GetBlog)
	r.POST("/blog", newBlogC.CreateBlog)
	r.DELETE("/blog/:id", newBlogC.DeleteBlog)
	r.PATCH("/blog/:id", newBlogC.UpdateBlog)

	r.Run(":8080")
}
