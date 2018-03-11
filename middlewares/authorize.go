package middlewares

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/stefanhipfel/lens-wide-shut/data"
)

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := appengine.NewContext(c.Request())
		u := user.Current(ctx)
		db, err := data.NewUserDatastoreDB(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		session, _ := session.Get("lws", c)
		v := session.Values["email"]
		if v != nil {
			return next(c)
		}

		if u == nil {
			url, _ := user.LoginURL(ctx, c.Request().URL.String())
			return c.Redirect(http.StatusTemporaryRedirect, url)

		}
		if user, err := db.GetUserByMail(u.Email); err == nil {
			if user.Email != "" {
				session.Values["email"] = u.Email
				session.Save(c.Request(), c.Response())
				return next(c)
			}
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "")
	}
}
