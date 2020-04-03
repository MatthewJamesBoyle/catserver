package cat_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/matthewjamesboyle/catserver/internal/cat"
	"github.com/matthewjamesboyle/catserver/internal/mock/mockcat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("Returns a Service and no error given valid input parameters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := mockcat.NewMockFactGetter(ctrl)
		g := mockcat.NewMockImageGetter(ctrl)

		s, err := cat.NewService(g, f)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	t.Run("Returns an error given a nil imageGetter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := mockcat.NewMockFactGetter(ctrl)
		s, err := cat.NewService(nil, f)

		assert.Nil(t, s)
		assert.Error(t, err)
		var e cat.ErrNilParam

		require.True(t, errors.As(err, &e))
		assert.Equal(t, "ImageGetter", e.Parameter)
	})

	t.Run("Returns an error given a nil FactGetter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := mockcat.NewMockImageGetter(ctrl)
		s, err := cat.NewService(g, nil)

		assert.Nil(t, s)
		assert.Error(t, err)
		var e cat.ErrNilParam

		require.True(t, errors.As(err, &e))
		assert.Equal(t, "FactGetter", e.Parameter)
	})
}

func TestService_GetImageAndFact(t *testing.T) {
	t.Run("given a valid request returns a catResult and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := mockcat.NewMockFactGetter(ctrl)
		g := mockcat.NewMockImageGetter(ctrl)
		s, err := cat.NewService(g, f)
		require.NoError(t, err)
		require.NotNil(t, s)

		someImage := cat.ImageURL("some-image-url")
		someFact := cat.Fact("some-fact")
		ctx := context.Background()
		ctxWc, _ := context.WithCancel(ctx)

		f.EXPECT().GetFact(ctxWc).Return(someFact, nil)
		g.EXPECT().GetImage(ctxWc).Return(someImage, nil)

		c, err := s.GetImageAndFact(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, someImage, c.ImageURL)
		assert.Equal(t, someFact, c.Fact)

	})

	t.Run("Returns an error given FactGetter Fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := mockcat.NewMockFactGetter(ctrl)
		g := mockcat.NewMockImageGetter(ctrl)
		s, err := cat.NewService(g, f)
		testErr := errors.New("some-error")
		ctx := context.Background()
		ctxWc, _ := context.WithCancel(ctx)

		f.EXPECT().GetFact(ctxWc).Return(cat.Fact(""), testErr)
		g.EXPECT().GetImage(ctxWc).Times(1)

		c, err := s.GetImageAndFact(context.Background())
		require.Equal(t, c, cat.CatResult{})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, testErr))
	})

	t.Run("Returns an error given FactGetter succeeds but imageGetter fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		f := mockcat.NewMockFactGetter(ctrl)
		g := mockcat.NewMockImageGetter(ctrl)
		s, err := cat.NewService(g, f)
		testErr := errors.New("some-error")
		ctx := context.Background()
		f.EXPECT().GetFact(gomock.Any()).Return(cat.Fact("some-fact"), nil)
		g.EXPECT().GetImage(gomock.Any()).Return(cat.ImageURL(""), testErr)

		c, err := s.GetImageAndFact(ctx)

		require.Equal(t, c, cat.CatResult{})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, testErr))

	})
}
