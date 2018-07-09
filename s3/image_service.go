package s3

import (
	"bytes"
	"image"
	_ "image/gif"  // import to register
	_ "image/jpeg" // import to register
	_ "image/png"  // import to register
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/j0hnsmith/progimage"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"

	"github.com/google/uuid"
)

var _ progimage.ImageService = &ImageService{}

// ImageService implements progimage.ImageService by storing data in S3 (or other compatible api).
type ImageService struct {
	BucketName string
	Client     *minio.Client
	UUID       func() uuid.UUID
}

// NewImageService provides an initialised ImageService.
func NewImageService(bucketName string, c *minio.Client, uuid func() uuid.UUID) *ImageService {
	return &ImageService{
		BucketName: bucketName,
		Client:     c,
		UUID:       uuid,
	}
}

// EnsureBucket creates the bucket (in us-east-1) if it doesn't already exist.
func (is *ImageService) EnsureBucket() error {
	exists, err := is.Client.BucketExists(is.BucketName)
	if err != nil {
		return errors.Wrap(err, "error checking bucket exists")
	}
	if exists {
		return nil
	}
	if err := is.Client.MakeBucket(is.BucketName, ""); err != nil {
		return errors.Wrap(err, "error creating bucket")
	}
	return nil
}

// Get retrieves the Image with the given id.
func (is *ImageService) Get(ID string) (progimage.Image, error) {
	ret := progimage.Image{}
	obj, err := is.Client.GetObject(is.BucketName, ID, minio.GetObjectOptions{})
	if err != nil {
		return ret, errors.Wrapf(err, "error getting image %s", ID)
	}

	// ensure the image exists
	info, err := obj.Stat()
	if err != nil {
		er, ok := err.(minio.ErrorResponse)
		if ok && er.Code == "NoSuchKey" {
			return ret, progimage.ErrImageNotFound
		}
		return ret, errors.Wrapf(err, "error getting image data %s", ID)
	}

	ret.ID = ID
	ret.Data = obj
	ret.ContentType = info.ContentType
	return ret, nil
}

// Store validates data is an image (read into memory), persists the image and returns the id.
func (is *ImageService) Store(rawImg io.Reader) (string, error) {
	// limit max size
	lr := io.LimitReader(rawImg, 20*1024*1024) // 20mb, refactor to config object so value can be set/modified

	// extract the mime type from the header
	b := make([]byte, 20)
	if _, err := lr.Read(b); err != io.EOF && err != nil {
		return "", errors.Wrap(err, "unable to read image data")
	}
	contentType := http.DetectContentType(b)
	if !strings.HasPrefix(contentType, "image") {
		// not an image, bail
		return "", progimage.ErrUnrecognisedImageType
	}

	// create 2 readers of image data, one is used to decode to ensure we have a valid image, the other is used
	// to upload to s3, both things happen at the same time. In the event that data is not a valid image, we
	// delete the uploaded object from s3. In theory we don't need to hold both the image bytes and the Image
	// object in memory so this should improve performance.

	// create 2 readers of rawImg (reads need to be syncronised, will block otherwise)
	pr, pw := io.Pipe()
	tr := io.TeeReader(io.MultiReader(bytes.NewReader(b), rawImg), pw)

	var err error
	errCh := make(chan error, 1)

	u := is.UUID()
	go func() {
		_, putErr := is.Client.PutObject(
			is.BucketName, u.String(),
			pr, -1,
			minio.PutObjectOptions{ContentType: contentType},
		)
		errCh <- putErr
	}()

	// decode the image to ensure we have a valid image
	var uploadErr error
	var decodeErr error
	if _, _, decodeErr = image.Decode(tr); decodeErr != nil {
		// io.EOF to read side
		pw.Close() // nolint: gas,errcheck
		uploadErr = <-errCh

		// delete uploaded image
		if err = is.Client.RemoveObject(is.BucketName, u.String()); err != nil {
			if uploadErr != nil {
				// Let's assume not uploaded to avoid further complexity in this example
				return "", progimage.ErrUnrecognisedImageType
			}

			// file not valid image, file uploaded ok but delete failed
			log.Printf("error deleting invalid image %s, %s", u, err)
			return "", progimage.ErrUnrecognisedImageType
		}

		return "", progimage.ErrUnrecognisedImageType
	}

	// nolint: gas,errcheck
	pw.Close() // io.EOF to read side
	uploadErr = <-errCh

	if uploadErr != nil {
		return "", errors.Wrap(uploadErr, "error uploading image to s3")
	}

	return u.String(), nil
}
