package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/stefanhipfel/lens-wide-shut/database"
	model "github.com/stefanhipfel/lens-wide-shut/models"

	"github.com/labstack/echo"
)

func Blog(c echo.Context) error {

	type blogData struct {
		Blog       *model.BlogPost
		Comments   []*model.Comment
		Like       *model.Like
		CoverPhoto *string
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	dbB, err := database.NewBlogDatastoreDB(c.Request())
	dbC, err := database.NewCommentDatastoreDB(c.Request())
	dbU, err := database.NewUserDatastoreDB(c.Request())
	dbL, err := database.NewLikeDatastoreDB(c.Request())
	dbI, err := database.NewImageStorage(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	blogPost, err := dbB.GetBlog(id)
	comments, err := dbC.GetCommentsByBlogId(id)
	user, err := dbU.GetUser(blogPost.CreatedBy)
	like, err := dbL.GetLike(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	blogPost.User = user

	cv, err := dbI.GetCoverPhotoFor(id, "blogs", 2048, false)

	if err != nil {
		cv = "../images/slider/slide-13-page-intro@2x.jpg"
	}

	d := blogData{blogPost, comments, like, &cv}

	var t *template.Template
	t, err = t.ParseFiles("./templates/blog-post.tmpl", "./templates/comment.tmpl")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = t.ExecuteTemplate(c.Response(), "blog", d)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return nil
}

func SaveBlog(c echo.Context) error {
	blogPost := new(model.BlogPost)
	err := c.Bind(blogPost)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db, err := database.NewBlogDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	id, err := db.AddBlog(blogPost)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

func GetBlogByUser(c echo.Context) error {
	user := c.QueryParam("user")

	if len(user) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing user")
	}

	db, err := database.NewBlogDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	blogPost, err := db.GetBlog(334)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, blogPost)
}

func GetBlogById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	type blogData struct {
		Blog     *model.BlogPost
		Comments []*model.Comment
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing id: "+err.Error())
	}

	db, err := database.NewBlogDatastoreDB(c.Request())
	dbC, err := database.NewCommentDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	blogPost, err := db.GetBlog(id)
	comments, err := dbC.GetCommentsByBlogId(id)

	d := blogData{blogPost, comments}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, d)
}

func GetAllBlogs(c echo.Context) error {
	db, err := database.NewBlogDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	blogPosts, err := db.ListBlogs()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, blogPosts)
}
