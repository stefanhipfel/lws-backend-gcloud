package common

import (
	"net/http"

	"github.com/labstack/echo"
)

func Startup(c echo.Context) error {
	//createDbSession()
	//CreateDatastoreClient(c.Request())

	return c.String(http.StatusOK, "")
}
