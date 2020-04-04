package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/matthewjamesboyle/catserver/internal/cat"
	"github.com/matthewjamesboyle/catserver/internal/mock/mockcat"
	"github.com/matthewjamesboyle/catserver/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHttpHandler(t *testing.T) {
	t.Run("returns an error given a nil servicer", func(t *testing.T) {
		h, err := transport.NewHttpHandler(nil)

		assert.Nil(t, h)
		assert.Error(t, err)
	})

	t.Run("Returns a httpHandler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mockcat.NewMockServicer(ctrl)
		h, err := transport.NewHttpHandler(s)

		assert.NoError(t, err)
		assert.IsType(t, &transport.HttpHandler{}, h)
	})
}

func TestHttpHandler_Get(t *testing.T) {
	t.Run("Returns a 500 given an error from servicer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		s := mockcat.NewMockServicer(ctrl)
		h, err := transport.NewHttpHandler(s)
		require.NoError(t, err)

		s.EXPECT().GetImageAndFact(ctx).Return(cat.CatResult{}, errors.New("some-error"))

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		h.Get(rr, r)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Returns a 200 and a catResult", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		someFact := cat.CatResult{
			ImageURL: "http//someurl",
			Fact:     "some-fact",
		}
		ctx := context.Background()
		s := mockcat.NewMockServicer(ctrl)
		h, err := transport.NewHttpHandler(s)
		require.NoError(t, err)

		s.EXPECT().GetImageAndFact(ctx).Return(someFact, nil)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		h.Get(rr, r)

		var res cat.CatResult

		err = json.NewDecoder(rr.Body).Decode(&res)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, someFact.ImageURL, res.ImageURL)
		assert.Equal(t, someFact.Fact, res.Fact)
	})

}
