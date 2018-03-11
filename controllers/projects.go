package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/stefanhipfel/lens-wide-shut/database"
)

func GetAllProjects(c echo.Context) error {
	dbP, err := database.NewProjectDatastoreDB(c.Request())
	dbI, err := database.NewImageStorage(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	projects, err := dbP.ListProjects()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, project := range projects {
		c, _ := dbI.GetCoverPhotoFor(project.ID, "projects", 300, false)
		project.CoverURL = c
	}

	return c.JSON(http.StatusOK, projects)
}

func GetProjectById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing id: "+err.Error())
	}
	dbP, err := database.NewProjectDatastoreDB(c.Request())
	dbI, err := database.NewImageStorage(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	project, err := dbP.GetProject(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	i, err := dbI.GetProjectImages(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	project.Images = &i

	return c.JSON(http.StatusOK, project)
}
