package cat_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/matthewjamesboyle/catserver/internal/cat"
	"github.com/matthewjamesboyle/catserver/internal/mock/mockcat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestNewFactService(t *testing.T) {
	t.Run("Returns a *FactService and no error", func(t *testing.T) {

		f, err := cat.NewFactService(http.DefaultClient, "some-baseurl")
		assert.NotNil(t, f)
		assert.NoError(t, err)
	})

	t.Run("Returns an error given a nil doer", func(t *testing.T) {
		f, err := cat.NewFactService(nil, "")

		assert.Nil(t, f)
		assert.Error(t, err)

		var e cat.ErrNilParam
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "hc", e.Parameter)
	})

	t.Run("Returns an error given an empty base url", func(t *testing.T) {
		f, err := cat.NewFactService(http.DefaultClient, "")

		assert.Nil(t, f)
		assert.Error(t, err)

		var e cat.ErrNilParam
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "baseUrl", e.Parameter)
	})
}

func TestFactService_GetFact(t *testing.T) {
	t.Run("Returns a Fact and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := mockcat.NewMockDoer(ctrl)
		someFact := cat.Fact("Jaguars are the only big cats that don't roar.")
		baseUrl := "http://some-baseurl"
		ctx := context.Background()
		u, err := url.Parse(baseUrl)
		require.NoError(t, err)

		sampleRes := `
		{
			"used":false,
			"source":"api",
			"type":"cat",
			"deleted":false,
			"_id":"591f98703b90f7150a19c180",
			"__v":0,
			"text":"Jaguars are the only big cats that don't roar.",
			"updatedAt":"2020-01-02T02:02:48.616Z",
			"createdAt":"2018-01-04T01:10:54.673Z",
			"status":{"verified":true,"sentCount":1},
			"user":"5a9ac18c7478810ea6c06381"
		}
		 `
		u.Path = "/facts"

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		require.NoError(t, err)

		s, err := cat.NewFactService(md, baseUrl)
		require.NoError(t, err)

		md.EXPECT().Do(req).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(sampleRes)),
		}, nil)

		f, err := s.GetFact(ctx)

		assert.NoError(t, err)
		assert.Equal(t, someFact, f)
	})
	t.Run("Returns an error given the http call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		d := mockcat.NewMockDoer(ctrl)

		ctx := context.Background()
		testErr := errors.New("some-error")
		baseURL := "http://someURL"
		u, err := url.Parse(baseURL)
		require.NoError(t, err)
		u.Path = "/facts"

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		require.NoError(t, err)

		d.EXPECT().Do(req).Return(nil, testErr)

		s, err := cat.NewFactService(d, baseURL)
		assert.NoError(t, err)

		f, err := s.GetFact(ctx)

		assert.Empty(t, f)
		assert.Error(t, err)

	})
}
