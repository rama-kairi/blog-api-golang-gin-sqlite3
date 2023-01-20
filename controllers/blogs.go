package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rama-kairi/blog-api-golang-gin/utils"
)

type Blog struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type BlogStore struct {
	Blogs []Blog
	db    *sql.DB
}

func NewBlogStore(dbIns *sql.DB) *BlogStore {
	return &BlogStore{
		Blogs: []Blog{},
		db:    dbIns,
	}
}

// Get all blogs
func (t BlogStore) GetAllBlogs(c *gin.Context) {
	rows, err := t.db.Query("SELECT * FROM blog")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error getting blogs")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var blog Blog
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Body)
		if err != nil {
			utils.Response(c, http.StatusInternalServerError, nil, "Error getting blogs")
			return
		}
		t.Blogs = append(t.Blogs, blog)
	}
	utils.Response(c, http.StatusOK, t.Blogs, "Blogs found")
}

// Get a blog
func (t BlogStore) GetBlog(c *gin.Context) {
	// Get blog id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Response(c, http.StatusBadRequest, nil, "Error getting blog")
		return
	}

	// Get the blog from the database
	row := t.db.QueryRow("SELECT * FROM blog WHERE id = ?", id)
	var blog Blog
	err = row.Scan(&blog.Id, &blog.Title, &blog.Body)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error getting blog")
		return
	}

	utils.Response(c, http.StatusNotFound, blog, "Blog found")
}

// Create a blog
func (t BlogStore) CreateBlog(c *gin.Context) {
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		utils.Response(c, http.StatusBadRequest, nil, "Error creating blog")
		return
	}

	// Save the blog to the database
	stmt, err := t.db.Prepare("INSERT INTO blog(title, body) VALUES(?, ?)")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error creating blog")
		return
	}
	res, err := stmt.Exec(blog.Title, blog.Body)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error creating blog")
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error creating blog")
		return
	}
	// Marshal the blog into json
	utils.Response(c, http.StatusCreated, id, "Blog created successfully")
}

// Delete a blog
func (t BlogStore) DeleteBlog(c *gin.Context) {
	// Get blog id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Response(c, http.StatusBadRequest, nil, "Error getting blog")
		return
	}
	// Delete the blog from the database
	stmt, err := t.db.Prepare("DELETE FROM blog WHERE id = ?")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error deleting blog")
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error deleting blog")
		return
	}

	// If the blog is not found, return 404
	utils.Response(c, http.StatusNoContent, nil, "Blog Deleted")
}

// Update a blog
func (t BlogStore) UpdateBlog(c *gin.Context) {
	// Get blog id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Response(c, http.StatusBadRequest, nil, "Error getting blog")
		return
	}

	// Get the blog from the request body
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		utils.Response(c, http.StatusBadRequest, nil, "Error updating blog")
		return
	}

	// Update the blog in the database
	stmt, err := t.db.Prepare("UPDATE blog SET title = ?, body = ? WHERE id = ?")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error updating blog")
		return
	}
	_, err = stmt.Exec(blog.Title, blog.Body, id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, nil, "Error updating blog")
		return
	}

	// If the blog is not found, return 404
	utils.Response(c, http.StatusNoContent, nil, "Blog Updated")
}
