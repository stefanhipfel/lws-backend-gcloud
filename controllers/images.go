package controllers

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/labstack/echo"

	"cloud.google.com/go/storage"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
)

type Image struct {
	Url string
}

type (
	Template struct {
		templates *template.Template
	}
)

func Images(c echo.Context) error {
	var images []Image
	ctx := appengine.NewContext(c.Request())

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "lens-wide-shut.appspot.com"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	query := &storage.Query{Prefix: "images/", Delimiter: "/"}
	it := bucket.Objects(ctx, query)

	s := c.QueryParam("size")
	imgSize, err := strconv.Atoi(s)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return echo.ErrNotFound
		}
		if !strings.Contains(attrs.Name, ".") {
			continue
		}
		filePath := attrs.Name
		blobPath := fmt.Sprintf("/gs/%s/%s", bucketName, filePath)
		bKey, err := blobstore.BlobKeyForFile(ctx, blobPath)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		opts := image.ServingURLOptions{Secure: false, Crop: true, Size: imgSize}
		url, err := image.ServingURL(ctx, bKey, &opts)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		image := Image{Url: url.String()}
		images = append(images, image)
	}
	t := template.New("image.html")                 // Create a template.
	t, err = t.ParseFiles("./templates/image.html") // Parse template file.
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	err = t.Execute(c.Response(), images) // merge.
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return nil
}
