package png

import (
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	primage "github.com/j0hnsmith/progimage/image"
)

var Transformer = primage.Transformer{
	Name:        "png",
	ContentType: "image/png",
	Encoder:     png.Encode,
}
