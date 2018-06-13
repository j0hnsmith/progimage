package progimage

import "io"

// Image represents a digital image.
type Image struct {
	ID          string
	Data        io.Reader
	ContentType string
}

// ImageService is an iterface for a service that can store and retrieve images.
type ImageService interface {
	Get(ID string) (Image, error)
	Store(imgRdr io.Reader) (string, error)
}

// ImageTypeTransformer is an interface that can transform images.
type ImageTypeTransformer interface {
	Transform(Image, chan error) (Image, error)
}
