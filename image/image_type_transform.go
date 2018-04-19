package imagetransform

import (
	"image"
	"io"

	"fmt"

	"bytes"

	"github.com/j0hnsmith/progimage"
	"github.com/pkg/errors"
)

var _ progimage.ImageTypeTransformer = Transformer{}

// Transformer enables progimage.ImageTypeTransformer implementations to be created easily avoiding code duplication.
type Transformer struct {
	ContentType string
	Encoder     func(io.Writer, image.Image) error
	Name        string
}

func (t Transformer) Transform(img progimage.Image) (progimage.Image, error) {
	if img.ContentType == t.ContentType {
		return img, nil
	}
	ret := progimage.Image{}
	i, _, err := image.Decode(img.Data)
	if err != nil {
		return ret, errors.Wrap(err, fmt.Sprintf("unable to decode %s image", t.Name))
	}
	b := new(bytes.Buffer) // using io.Pipe() causes a deadlock (obviously), would be nice not to write into memory
	if err := t.Encoder(b, i); err != nil {
		return ret, errors.Wrap(err, fmt.Sprintf("unable to encode %s image", t.Name))
	}
	ret.ID = img.ID
	ret.ContentType = t.ContentType
	ret.Data = b
	return ret, nil
}
