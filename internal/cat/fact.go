package cat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Fact string

type FactGetter interface {
	GetFact(ctx context.Context) (Fact, error)
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type FactService struct {
	hc      Doer
	baseUrl string
}

type FactResponse struct {
	Text string `json:"text"`
}

func NewFactService(hc Doer, baseUrl string) (*FactService, error) {

	if hc == nil {
		return nil, ErrNilParam{Parameter: "hc"}
	}

	if baseUrl == "" {
		return nil, ErrNilParam{Parameter: "baseUrl"}
	}

	return &FactService{
		hc:      hc,
		baseUrl: baseUrl,
	}, nil
}

func (f *FactService) GetFact(ctx context.Context) (Fact, error) {
	u, err := url.Parse(f.baseUrl)
	if err != nil {
		return "", err
	}
	u.Path = "/facts"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	resp, err := f.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling fact service: %w", err)
	}

	var fr FactResponse
	err = json.NewDecoder(resp.Body).Decode(&fr)
	if err != nil {
		return "", fmt.Errorf("unmarshall Response: %w", err)
	}

	return Fact(fr.Text), nil
}
