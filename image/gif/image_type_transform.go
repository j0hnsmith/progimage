package gif

import (
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	primage "github.com/j0hnsmith/progimage/image"
)

var Transformer = primage.Transformer{
	Name:        "gif",
	ContentType: "image/gif",
	Encoder:     DefaultGifEncode,
}

// DefaultGifEncode performs fig encoding with default values for options.
func DefaultGifEncode(w io.Writer, m image.Image) error {
	return gif.Encode(w, m, nil)
}
