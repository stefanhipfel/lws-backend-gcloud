package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/stefanhipfel/lens-wide-shut/database"
)

func Index(c echo.Context) error {
	var err error

	dbBlog, err := database.NewBlogDatastoreDB(c.Request())
	dbUser, err := database.NewUserDatastoreDB(c.Request())
	dbI, err := database.NewImageStorage(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	blogPosts, err := dbBlog.ListBlogs()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, k := range blogPosts {
		user, err := dbUser.GetUser(blogPosts[i].CreatedBy)
		if err != nil {
			return fmt.Errorf("datastoredb: could not list blogs: %v", err)
		}
		blogPosts[i].User = user
		cf, err := dbI.GetCoverPhotoFor(k.ID, "blogs", 600, false)

		if err != nil {
			cf = "images/photo-studio/blog/image-regular-1.jpg"
		}
		blogPosts[i].ImageUrl = cf

	}

	t := template.New("index.html")                 // Create a template.
	t, err = t.ParseFiles("./templates/index.html") // Parse template file.
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	err = t.Execute(c.Response(), blogPosts) // merge.
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return nil
}
