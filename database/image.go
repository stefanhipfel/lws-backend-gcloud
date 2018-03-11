package database

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/stefanhipfel/lens-wide-shut/models"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	aImage "google.golang.org/appengine/image"
)

const bucketName = "lens-wide-shut.appspot.com"

func NewImageStorage(r *http.Request) (*StorageDB, error) {
	var ctx context.Context
	ctx = appengine.NewContext(r)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	return &StorageDB{
		bucket: bucket,
		ctx:    ctx,
	}, nil
}

func (db *StorageDB) GetCoverPhotoFor(id int64, folder string, size int, crop bool) (string, error) {
	obj := db.bucket.Object(folder + "/" + strconv.FormatInt(id, 10) + "/cover.jpg")
	i := models.Image{}
	attrs, err := obj.Attrs(db.ctx)
	if err != nil {
		return "", fmt.Errorf("image: could not get object attrs: %v", err)
	}
	err = db.setImageURL(&i, attrs.Name, size, crop)
	if err != nil {
		return "", err
	}
	return i.Url, nil
}

func (db *StorageDB) GetProjectImages(id int64) ([]models.Image, error) {
	query := &storage.Query{Prefix: "projects/" + strconv.FormatInt(id, 10) + "/", Delimiter: "/"}
	it := db.bucket.Objects(db.ctx, query)
	imges := make([]models.Image, 0)
	c := make(chan models.Image)
	var wg sync.WaitGroup
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return imges, err
		}
		wg.Add(1)
		go func(attrs storage.ObjectAttrs) {
			img := models.Image{}
			defer wg.Done()
			if !strings.Contains(attrs.Name, ".") {
				return
			}

			err := db.setImageURL(&img, attrs.Name, 600, false)
			if err != nil {
				return
			}

			err = db.setImageSize(&img, attrs.Name)
			if err != nil {
				return
			}
			c <- img
		}(*attrs)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for img := range c {
		imges = append(imges, img)
	}

	return imges, nil
}

func (db *StorageDB) setImageURL(i *models.Image, filePath string, size int, crop bool) error {
	var url *url.URL
	blobPath := fmt.Sprintf("/gs/%s/%s", bucketName, filePath)
	bKey, err := blobstore.BlobKeyForFile(db.ctx, blobPath)
	if err != nil {
		return fmt.Errorf("image: could not get blobkey for image: %v", err)
	}
	opts := aImage.ServingURLOptions{Secure: false, Crop: crop, Size: size}
	url, err = aImage.ServingURL(db.ctx, bKey, &opts)
	if err != nil {
		return fmt.Errorf("image: could not get url for image: %v", err)
	}

	i.Url = url.String()

	return nil
}

func (db *StorageDB) setImageSize(i *models.Image, name string) error {
	r, err := db.bucket.Object(name).NewReader(db.ctx)

	if err != nil {
		return fmt.Errorf("image: could create object reader: %v", err)
	}
	defer r.Close()

	c, _, err := image.DecodeConfig(r)

	if err != nil {
		return fmt.Errorf("image: could not decode image: %v", err)
	}

	mx := math.Max(float64(c.Width), float64(c.Height))
	mi := math.Min(float64(c.Width), float64(c.Height))
	rt := mx / mi

	i.Width = c.Width
	i.Height = c.Height
	i.Ratio = rt

	return nil
}
