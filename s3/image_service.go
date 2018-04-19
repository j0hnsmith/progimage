package s3

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"

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
	// validate we have an image that we can process
	// optimisation: to avoid reading entire image into memory, https://golang.org/pkg/net/http/#DetectContentTypec,
	// doesn't mean that an image is valid though
	b := &bytes.Buffer{}
	tr := io.TeeReader(rawImg, b)

	var typ string
	var err error
	if _, typ, err = image.Decode(tr); err != nil {
		return "", progimage.ErrUnrecognisedImageType
	}

	u := is.UUID()
	_, err = is.Client.PutObject(
		is.BucketName, u.String(),
		b, int64(b.Len()),
		minio.PutObjectOptions{ContentType: mime.TypeByExtension("." + typ)},
	)
	if err != nil {
		return "", errors.Wrap(err, "error uploading image to s3")
	}

	return u.String(), nil
}
