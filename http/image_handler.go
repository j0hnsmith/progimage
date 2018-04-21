package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/j0hnsmith/progimage"
	"github.com/j0hnsmith/progimage/image/gif"
	"github.com/j0hnsmith/progimage/image/jpeg"
	"github.com/j0hnsmith/progimage/image/png"
	"github.com/julienschmidt/httprouter"
)

const MaxReadBytes = 50 * 1024 * 1024 // 50mb

type ImageHandler struct {
	*httprouter.Router

	Transformers map[string]progimage.ImageTypeTransformer
	ImageService progimage.ImageService
}

// NewImageHandler returns an initialised image handler.
func NewImageHandler(is progimage.ImageService) *ImageHandler {
	h := ImageHandler{
		Router:       httprouter.New(),
		ImageService: is,
		Transformers: map[string]progimage.ImageTypeTransformer{
			"png": png.Transformer,
			"jpg": jpeg.Transformer,
			"gif": gif.Transformer,
		},
	}
	h.POST("/image/create", h.handleCreateImage)
	h.GET("/image/:id", h.handleGetImage)
	return &h
}

func (h ImageHandler) handleCreateImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// don't allow an attacker to send an unlimited stream of bytes
	lr := io.LimitReader(r.Body, MaxReadBytes)

	ID, err := h.ImageService.Store(lr)
	if err != nil {
		if err == progimage.ErrUnrecognisedImageType {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": "%s"}`, ID)))
}

func (h ImageHandler) handleGetImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ID := params.ByName("id")

	s := strings.Split(ID, ".")
	if len(s) == 2 {
		h.handleGetImageWithExt(w, r, s[0], s[1])
		return
	}

	h.handleGetImageNoExt(w, r, ID)
	return
}

func (h ImageHandler) handleGetImageNoExt(w http.ResponseWriter, r *http.Request, ID string) {
	img, err := h.ImageService.Get(ID)
	if err != nil {
		if err == progimage.ErrImageNotFound {
			http.Error(w, fmt.Sprintf("image %s not found", ID), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", img.ContentType)
	io.Copy(w, img.Data)
}

func (h ImageHandler) handleGetImageWithExt(w http.ResponseWriter, r *http.Request, ID, ext string) {
	tr, ok := h.Transformers[ext]
	if !ok {
		http.Error(w, "unsupported image type", http.StatusBadRequest)
		return
	}
	imgOrig, err := h.ImageService.Get(ID)
	if err != nil {
		if err == progimage.ErrImageNotFound {
			http.Error(w, fmt.Sprintf("image %s not found", ID), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ec := make(chan error, 1)
	imgConv, err := tr.Transform(imgOrig, ec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", imgConv.ContentType)
	written, err := io.Copy(w, imgConv.Data)
	if err != nil {
		if written == 0 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// 200 sent already, all we can do is log
			log.Printf(
				"error converting %s to %s (id: %s), 200 sent already",
				imgOrig.ContentType,
				imgConv.ContentType,
				imgOrig.ID,
			)
		}
	}
}
