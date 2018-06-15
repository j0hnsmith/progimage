package png

import (
	_ "image/gif"  // register image type, do not remove
	_ "image/jpeg" // register image type, do not remove
	"image/png"

	primage "github.com/j0hnsmith/progimage/image"
)

// Transformer implements progimage.ImageTypeTransformer to convert a progimage.Image to png format.
var Transformer = primage.Transformer{
	Name:        "png",
	ContentType: "image/png",
	Encoder:     png.Encode,
}
