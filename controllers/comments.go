package controllers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/stefanhipfel/lens-wide-shut/database"
	"github.com/stefanhipfel/lens-wide-shut/models"
)

func GetCommentsByBlogId(c echo.Context) error {
	return nil
}

func SaveComment(c echo.Context) error {
	bId, err := strconv.ParseInt(c.FormValue("blogId"), 10, 64)
	db, err := database.NewCommentDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fn := c.FormValue("fname")
	w := c.FormValue("website")
	ms := c.FormValue("message")
	em := c.FormValue("email")

	if len(fn) == 0 || len(ms) == 0 || len(em) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing key fields")
	}

	co := models.Comment{
		Name:      fn,
		Website:   w,
		Message:   ms,
		BlogID:    bId,
		CreatedBy: em,
		CreatedOn: time.Now(),
	}

	id, err := db.AddComment(&co)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

func GetComment(c echo.Context) error {
	type d struct {
		Comments []*models.Comment
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db, err := database.NewCommentDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	co, err := db.GetCommentsByBlogId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var t *template.Template
	t, err = t.ParseFiles("./templates/comment.tmpl")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = t.ExecuteTemplate(c.Response(), "comment", d{Comments: co})
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return nil

	return c.JSON(http.StatusOK, co)

}
