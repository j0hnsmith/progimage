package mock

import (
	"io"

	"github.com/j0hnsmith/progimage"
)

var _ progimage.ImageService = &ImageService{}

// ImageService is a mock progimage.ImageService.
type ImageService struct {
	GetInvoked   bool
	StoreInvoked bool
	GetFunc      func(string) (progimage.Image, error)
	StoreFunc    func(io.Reader) (string, error)
}

// Get an image.
func (is *ImageService) Get(ID string) (progimage.Image, error) {
	is.GetInvoked = true
	return is.GetFunc(ID)
}

// Store an image.
func (is *ImageService) Store(imgRdr io.Reader) (string, error) {
	is.StoreInvoked = true
	return is.StoreFunc(imgRdr)
}
