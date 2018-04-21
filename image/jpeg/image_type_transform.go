package jpeg

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"

	primage "github.com/j0hnsmith/progimage/image"
)

// Transformer implements progimage.ImageTypeTransformer to convert a progimage.Image to jpeg format.
var Transformer = primage.Transformer{
	Name:        "jpeg",
	ContentType: "image/jpeg",
	Encoder:     DefaultJpegEncode,
}

// DefaultJpegEncode performs jpeg encoding with default values for options.
func DefaultJpegEncode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}
