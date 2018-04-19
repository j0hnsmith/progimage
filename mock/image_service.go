package mock

import (
	"io"

	"github.com/j0hnsmith/progimage"
)

var _ progimage.ImageService = &ImageService{}

type ImageService struct {
	GetInvoked   bool
	GetFunc      func(string) (progimage.Image, error)
	StoreInvoked bool
	StoreFunc    func(io.Reader) (string, error)
}

func (is *ImageService) Get(ID string) (progimage.Image, error) {
	is.GetInvoked = true
	return is.GetFunc(ID)
}

func (is *ImageService) Store(imgRdr io.Reader) (string, error) {
	is.StoreInvoked = true
	return is.StoreFunc(imgRdr)
}
