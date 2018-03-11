// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package static demonstrates a static file handler for App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/controllers"
	"github.com/stefanhipfel/lens-wide-shut/middlewares"
	"github.com/stefanhipfel/notsureyet/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	_ "github.com/valyala/bytebufferpool"
	"google.golang.org/appengine"
)

func main() {
	// the appengine package provides a convenient method to handle the health-check requests
	// and also run the app on the correct port. We just need to add Echo to the default handler
	e := echo.New()
	http.Handle("/", e)

	store := sessions.NewCookieStore([]byte(handlers.RandomToken(64)))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	e.Use(session.Middleware(sessions.Store(store)))

	e.GET("/_ah/start", common.Startup)

	e.GET("/", controllers.Index)
	e.GET("/images", controllers.Images)
	e.GET("/blog/:id", controllers.Blog)
	e.GET("/blogs/:id", controllers.GetBlogById)
	e.PUT("/like/:id", controllers.Vote)
	e.GET("/like/:id", controllers.GetLike)
	e.POST("/comments", controllers.SaveComment)
	e.GET("/comments/:id", controllers.GetComment)
	e.POST("/contact", controllers.Contact)

	e.GET("/api/projects", controllers.GetAllProjects)
	e.GET("/api/projects/:id", controllers.GetProjectById)

	authorized := e.Group("/")
	authorized.Use(middlewares.AuthorizeRequest)

	authorized.GET("/blogs/:id", controllers.GetBlogById)
	authorized.GET("/user/:email", controllers.GetUserByMail)
	authorized.POST("/blogs", controllers.SaveBlog)
	authorized.GET("/blogs", controllers.GetAllBlogs)
	authorized.GET("/blogs:id/like", controllers.GetAllBlogs)
	//authorized.GET("/comments/:id", controllers.GetCommentsByBlogId)
	defer common.GetDatastoreClient().Close()

	appengine.Main()
}

// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
