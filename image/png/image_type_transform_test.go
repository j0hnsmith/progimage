package png_test

import (
	"image"
	"os"
	"testing"

	"github.com/j0hnsmith/progimage"
	"github.com/j0hnsmith/progimage/image/png"
)

var fileTests = []struct {
	Name        string
	Path        string
	ContentType string
}{
	{Name: "png", Path: "../../testimages/test.png", ContentType: "image/png"},
	{Name: "gif", Path: "../../testimages/test.gif", ContentType: "image/gif"},
	{Name: "jpg", Path: "../../testimages/test.jpg", ContentType: "image/jpeg"},
}

func TestTransformPNG(t *testing.T) {
	for _, item := range fileTests {
		t.Run(item.Name, func(t *testing.T) {
			fp, err := os.Open(item.Path)
			if err != nil {
				t.Fatal(err)
			}
			defer fp.Close()

			img := progimage.Image{
				ID:          item.Name,
				ContentType: item.ContentType,
				Data:        fp,
			}

			imgOut, err := png.Transformer.Transform(img)
			if err != nil {
				t.Fatal(err)
			}

			_, typ, err := image.Decode(imgOut.Data)
			if err != nil {
				t.Fatal(err)
			}
			if typ != "png" {
				t.Errorf("expected type of converted image to be png, got %s", typ)
			}
		})
	}
}