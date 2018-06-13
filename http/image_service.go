package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/j0hnsmith/progimage"
	"github.com/pkg/errors"
)

type ImageService struct {
	BaseUrl string
	Client  *http.Client
}

var _ progimage.ImageService = ImageService{}

func (is ImageService) Get(ID string) (progimage.Image, error) {
	resp, err := is.Client.Get(is.BaseUrl + "/image/" + ID)
	ret := progimage.Image{}
	if err != nil {
		return ret, errors.Wrap(err, "unable to make get request")
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return ret, progimage.ErrImageNotFound
		}
		return ret, errors.New(fmt.Sprintf("unknown error getting image, status code %s", resp.StatusCode))
	}

	ret.ID = ID
	ret.ContentType = resp.Header.Get("ContentType")
	ret.Data = resp.Body
	return ret, nil
}

func (is ImageService) Store(imgRdr io.Reader) (string, error) {
	req, err := http.NewRequest("POST", is.BaseUrl+"/image/create", imgRdr)
	if err != nil {
		return "", errors.Wrap(err, "unable to create new http request")
	}
	resp, err := is.Client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to make post request")
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.New(fmt.Sprintf("unknown error creating new image, status code %s", resp.StatusCode))
	}

	rd := new(respData)
	if err := json.NewDecoder(resp.Body).Decode(rd); err != nil {
		return "", errors.Wrap(err, "error decoding resp")
	}

	return rd.ID, nil
}

type respData struct {
	ID string
}
