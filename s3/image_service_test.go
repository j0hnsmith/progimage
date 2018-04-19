// package s3_test tests the s3 image service against a real s3 api implementation. The following env vars are required
// to run the tests
//
//     S3_ENDPOINT
//     S3_ACCESS_KEY
//     S3_SECRET_KEY
//     S3_ACCESS_KEY
//     S3_SECURE
//
// You can use real s3 creds, alternatively
// `docker run -p 9000:9000 -e MINIO_ACCESS_KEY=minio -e MINIO_SECRET_KEY=miniostorage minio/minio server /data` will
// provide a compatible service using minio.
package s3_test

import (
	"testing"

	"io/ioutil"

	"bytes"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/j0hnsmith/progimage"
	"github.com/j0hnsmith/progimage/s3"
	"github.com/minio/minio-go"
)

const testBucketName = "test-bucket-s3"

var c *minio.Client

func init() {
	// setup
	var endpoint, accessKey, secretKey string
	var secure bool

	if endpoint = os.Getenv("S3_ENDPOINT"); endpoint == "" {
		panic("env var S3_ENDPOINT is required")
	}
	if accessKey = os.Getenv("S3_ACCESS_KEY"); accessKey == "" {
		panic("env var S3_ACCESS_KEY is required")
	}
	if secretKey = os.Getenv("S3_SECRET_KEY"); secretKey == "" {
		panic("env var S3_SECRET_KEY is required")
	}
	if s := os.Getenv("S3_SECURE"); s == "" {
		panic("env var S3_SECURE is required")
	} else {
		secure = s == "true"
	}

	var err error
	c, err = minio.New(endpoint, accessKey, secretKey, secure)
	if err != nil {
		panic(err)
	}
}

// setup deletes all objects in the test bucket then deletes the bucket.
func setup(t *testing.T) {
	if exists, err := c.BucketExists(testBucketName); err != nil {
		t.Fatal(err)
	} else if !exists {
		// bucket doesn't exist, bail
		return
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	// must empty bucket before deleting it
	for obj := range c.ListObjects(testBucketName, "", true, doneCh) {
		err := c.RemoveObject(testBucketName, obj.Key)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := c.RemoveBucket(testBucketName); err != nil {
		t.Fatal(err)
	}
}

var fileTests = []struct {
	Name        string
	Path        string
	ContentType string
}{
	{Name: "png", Path: "../testimages/test.png", ContentType: "image/png"},
	{Name: "gif", Path: "../testimages/test.gif", ContentType: "image/gif"},
	{Name: "jpg", Path: "../testimages/test.jpg", ContentType: "image/jpeg"},
}

// Test storing and retrieving images of each type, a missed import will cause a failure, see
// https://golang.org/pkg/image/#pkg-overview
func TestImageService_StoreGetImage(t *testing.T) {
	for _, item := range fileTests {
		t.Run(item.Name, func(t *testing.T) {
			setup(t)

			uid := uuid.New()
			uf := func() uuid.UUID {
				return uid
			}

			is := s3.NewImageService(testBucketName, c, uf)
			if err := is.EnsureBucket(); err != nil {
				t.Fatal(err)
			}

			// store image
			fp, err := os.Open(item.Path)
			if err != nil {
				t.Fatal(err)
			}
			d, err := ioutil.ReadAll(fp)
			if err != nil {
				t.Fatal(err)
			}
			fp.Close()

			// keep a copy of initial data to compare to returned data
			initial := bytes.NewReader(d)

			// store then retrieve using id
			id, err := is.Store(initial)
			if err != nil {
				t.Fatal(err)
			}
			if id == "" {
				t.Error("expected id to be populated, got empty string")
			}
			if id != uid.String() {
				t.Errorf("expected id to be %s, got %s", uid.String(), id)
			}

			img, err := is.Get(id)
			if err != nil {
				t.Fatal(err)
			}

			if img.ID != id {
				t.Errorf("expected image id to be '%s', got '%s'", id, img.ID)
			}

			if img.ContentType != item.ContentType {
				t.Errorf("expected image content type to be '%s', got '%s'", item.ContentType, img.ContentType)
			}

			retrieved, err := ioutil.ReadAll(img.Data)
			if err != nil {
				t.Fatal(err)
			}

			b := new(bytes.Buffer)
			initial.Seek(0, io.SeekStart)
			_, err = b.ReadFrom(initial)
			if err != nil {
				t.Fatal(err)
			}

			// compare initial data to retrieved data
			if !bytes.Equal(b.Bytes(), retrieved) {
				t.Error("expected stored data to be equal to initial data from file, data not equal")
			}
		})
	}
}

func TestImageService_StoreNoData(t *testing.T) {
	setup(t)

	is := s3.NewImageService(testBucketName, c, uuid.New)
	if err := is.EnsureBucket(); err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader([]byte{})
	if _, err := is.Store(r); err != progimage.ErrUnrecognisedImageType {
		t.Errorf("expected progimage.ErrUnrecognisedImageType, got %s", err)
	}
}

func TestImageService_GetNotExists(t *testing.T) {
	setup(t)

	is := s3.NewImageService(testBucketName, c, uuid.New)
	if err := is.EnsureBucket(); err != nil {
		t.Fatal(err)
	}

	if _, err := is.Get("foo"); err != progimage.ErrImageNotFound {
		t.Errorf("expected progimage.ErrImageNotFound, got %s", err)
	}
}
