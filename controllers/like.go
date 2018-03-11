package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/stefanhipfel/lens-wide-shut/database"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func Vote(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := appengine.NewContext(c.Request())
	u := user.Current(ctx)
	if u == nil {
		u := strings.Replace(c.Request().URL.String(), "like", "blog", -1)
		url, _ := user.LoginURL(ctx, u)
		return c.String(http.StatusTemporaryRedirect, url)

	}
	db, err := database.NewLikeDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = db.Vote(id, u.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func GetLike(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	db, err := database.NewLikeDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	like, err := db.GetLike(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	l := len(like.Likes)
	return c.String(http.StatusOK, strconv.Itoa(l))
}
