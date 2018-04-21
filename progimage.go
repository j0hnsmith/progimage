package progimage

import "io"

type Image struct {
	ID          string
	Data        io.Reader
	ContentType string
}

type ImageService interface {
	Get(ID string) (Image, error)
	Store(imgRdr io.Reader) (string, error)
}

type ImageTypeTransformer interface {
	Transform(Image, chan error) (Image, error)
}
