package png

import (
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	primage "github.com/j0hnsmith/progimage/image"
)

// Transformer implements progimage.ImageTypeTransformer to convert a progimage.Image to png format.
var Transformer = primage.Transformer{
	Name:        "png",
	ContentType: "image/png",
	Encoder:     png.Encode,
}
