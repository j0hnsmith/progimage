package http_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	pihttp "github.com/j0hnsmith/progimage/http"
	"github.com/j0hnsmith/progimage/mock"
)

func TestServer_StartStop(t *testing.T) {
	// simple start/stop test

	is := new(mock.ImageService)
	h := pihttp.NewImageHandler(is)

	s := pihttp.Server{
		ImageHandler: *h,
		Addr:         "127.0.0.1:34567",
	}

	var startErr error
	go func() {
		startErr = s.Start(new(bytes.Buffer))
	}()

	// arbitrary sleep to wait for server to start before stopping
	time.Sleep(time.Second * 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var stopErr error
	stopErr = s.Stop(ctx)

	if startErr != nil {
		t.Errorf("got error when starting (or stopping) server, %+v", startErr)
	}

	if stopErr != nil {
		t.Errorf("got error when stopping server, %+v", startErr)
	}
}
