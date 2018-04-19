package progimage

import "errors"

// ErrImageNotFound represents an image not found.
var ErrImageNotFound = errors.New("image not found")

// ErrUnrecognisedImageType represents image data that can't be processed.
var ErrUnrecognisedImageType = errors.New("unrecognised image data")
