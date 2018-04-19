package http

import (
	"net/http"

	"fmt"
	"io"

	"github.com/j0hnsmith/progimage"
	"github.com/julienschmidt/httprouter"
)

const MaxReadBytes = 50 * 1024 * 1024 // 50mb

type ImageHandler struct {
	*httprouter.Router

	ImageService progimage.ImageService
}

// NewImageHandler returns an initialised image handler.
func NewImageHandler(is progimage.ImageService) *ImageHandler {
	h := ImageHandler{
		Router:       httprouter.New(),
		ImageService: is,
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

	w.Header().Set("ContentType", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": "%s"}`, ID)))
}

func (h ImageHandler) handleGetImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ID := params.ByName("id")
	img, err := h.ImageService.Get(ID)
	if err != nil {
		if err == progimage.ErrImageNotFound {
			http.Error(w, fmt.Sprintf("image %s not found", ID), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("ContentType", img.ContentType)
	io.Copy(w, img.Data)
}
