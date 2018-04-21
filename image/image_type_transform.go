package imagetransform

import (
	"fmt"
	"image"
	"io"

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

func (t Transformer) Transform(img progimage.Image, ec chan error) (progimage.Image, error) {
	if img.ContentType == t.ContentType {
		ec <- nil
		return img, nil
	}
	ret := progimage.Image{}
	i, _, err := image.Decode(img.Data)
	if err != nil {
		return ret, errors.Wrap(err, fmt.Sprintf("unable to decode %s image", t.Name))
	}

	r, w := io.Pipe()
	go func() {
		if err := t.Encoder(w, i); err != nil {
			ec <- errors.Wrap(err, fmt.Sprintf("unable to encode %s image", t.Name))
			w.Close()
			return
		}
		ec <- nil
		w.Close()
	}()

	ret.ID = img.ID
	ret.ContentType = t.ContentType
	ret.Data = r
	return ret, nil
}
