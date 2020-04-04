package cat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ImageResponse []struct {
	URL string `json:"url"`
}

type ImageURL string

type ImageGetter interface {
	GetImage(ctx context.Context) (ImageURL, error)
}

type ImageService struct {
	url string
	hc  Doer
}

func NewImageService(hc Doer, url string) (*ImageService, error) {
	if hc == nil {
		return nil, ErrNilParam{Parameter: "Doer"}
	}

	if url == "" {
		return nil, ErrNilParam{Parameter: "url"}
	}

	return &ImageService{
		url: url,
		hc:  hc,
	}, nil
}

func (s *ImageService) GetImage(ctx context.Context) (ImageURL, error) {

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return "", err
	}

	res, err := s.hc.Do(r)
	if err != nil {
		return "", err
	}
	var x ImageResponse
	err = json.NewDecoder(res.Body).Decode(&x)
	if err != nil {
		return "", err
	}
	if len(x) == 0 {
		return "", errors.New("not long enough mate")
	}

	return ImageURL(x[0].URL), nil
}
