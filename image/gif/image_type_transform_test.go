package gif_test

import (
	"image"
	"os"
	"testing"

	"github.com/j0hnsmith/progimage"
	"github.com/j0hnsmith/progimage/image/gif"
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

func TestTransformGif(t *testing.T) {
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

			errCh := make(chan error, 1)
			imgOut, err := gif.Transformer.Transform(img, errCh)
			if err != nil {
				t.Fatal(err)
			}

			_, typ, err := image.Decode(imgOut.Data)
			if err != nil {
				t.Fatal(err)
			}
			if typ != "gif" {
				t.Errorf("expected type of converted image to be gif, got %s", typ)
			}
			if err := <-errCh; err != nil {
				t.Errorf("got error converting image %s", err)
			}
		})
	}
}
