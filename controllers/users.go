package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/stefanhipfel/lens-wide-shut/database"
	model "github.com/stefanhipfel/lens-wide-shut/models"
)

func GetUserByMail(c echo.Context) error {
	mail := c.Param("email")
	db, err := database.NewUserDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	blogPost, err := db.GetUserByMail(mail)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, blogPost)
}

func SaveUser(c echo.Context) error {
	user := new(model.User)

	err := c.Bind(user)
	db, err := database.NewUserDatastoreDB(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := db.AddUser(user)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
