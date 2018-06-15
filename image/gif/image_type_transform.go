package gif

import (
	"image"
	"image/gif"
	_ "image/jpeg" // import to register image type
	_ "image/png"  // import to register image type
	"io"

	primage "github.com/j0hnsmith/progimage/image"
)

// Transformer implements progimage.ImageTypeTransformer to convert a progimage.Image to png format.
var Transformer = primage.Transformer{
	Name:        "gif",
	ContentType: "image/gif",
	Encoder:     DefaultGifEncode,
}

// DefaultGifEncode performs fig encoding with default values for options.
func DefaultGifEncode(w io.Writer, m image.Image) error {
	return gif.Encode(w, m, nil)
}
