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
	"testing"
)

func TestNewImageService(t *testing.T) {
	t.Run("returns a ImageService and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		d := mockcat.NewMockDoer(ctrl)

		s, err := cat.NewImageService(d, "some-url")

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	t.Run("Returns an error given a nil doer", func(t *testing.T) {
		_, err := cat.NewImageService(nil, "")

		assert.Error(t, err)
		var e cat.ErrNilParam
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "Doer", e.Parameter)
	})

	t.Run("Returns an error given a nil image url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		d := mockcat.NewMockDoer(ctrl)
		_, err := cat.NewImageService(d, "")

		assert.Error(t, err)
		var e cat.ErrNilParam
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "url", e.Parameter)
	})
}

func TestImageService_GetImage(t *testing.T) {
	t.Run("returns an error if doer returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := "someurl"
		ctx := context.Background()
		d := mockcat.NewMockDoer(ctrl)
		s, err := cat.NewImageService(d, u)
		require.NoError(t, err)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		require.NoError(t, err)

		//		sampleRes := `
		//		[
		//   {
		//      "breeds":[
		//         {
		//            "weight":{
		//               "imperial":"7 - 14",
		//               "metric":"3 - 6"
		//            },
		//            "id":"esho",
		//            "name":"Exotic Shorthair",
		//            "cfa_url":"http://cfa.org/Breeds/BreedsCJ/Exotic.aspx",
		//            "vetstreet_url":"http://www.vetstreet.com/cats/exotic-shorthair",
		//            "vcahospitals_url":"https://vcahospitals.com/know-your-pet/cat-breeds/exotic-shorthair",
		//            "temperament":"Affectionate, Sweet, Loyal, Quiet, Peaceful",
		//            "origin":"United States",
		//            "country_codes":"US",
		//            "country_code":"US",
		//            "description":"The Exotic Shorthair is a gentle friendly cat that has the same personality as the Persian. They love having fun, don’t mind the company of other cats and dogs, also love to curl up for a sleep in a safe place. Exotics love their own people, but around strangers they are cautious at first. Given time, they usually warm up to visitors.",
		//            "life_span":"12 - 15",
		//            "indoor":0,
		//            "lap":1,
		//            "alt_names":"Exotic",
		//            "adaptability":5,
		//            "affection_level":5,
		//            "child_friendly":3,
		//            "dog_friendly":3,
		//            "energy_level":3,
		//            "grooming":2,
		//            "health_issues":3,
		//            "intelligence":3,
		//            "shedding_level":2,
		//            "social_needs":4,
		//            "stranger_friendly":2,
		//            "vocalisation":1,
		//            "experimental":0,
		//            "hairless":0,
		//            "natural":0,
		//            "rare":0,
		//            "rex":0,
		//            "suppressed_tail":0,
		//            "short_legs":0,
		//            "wikipedia_url":"https://en.wikipedia.org/wiki/Exotic_Shorthair",
		//            "hypoallergenic":0
		//         }
		//      ],
		//      "id":"y61B6bFCh",
		//      "url":"https://cdn2.thecatapi.com/images/y61B6bFCh.jpg",
		//      "width":898,
		//      "height":900
		//   }
		//]
		//		`

		//res := http.Response{
		//	Body: ioutil.NopCloser(bytes.NewBuffer(sampleRes)),
		//}
		d.EXPECT().Do(r).Return(nil, errors.New("some-error"))
		i, err := s.GetImage(ctx)

		assert.Empty(t, i)
		assert.Error(t, err)
	})

	t.Run("returns an error as it cannot decode the response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := "someurl"
		ctx := context.Background()
		d := mockcat.NewMockDoer(ctrl)
		s, err := cat.NewImageService(d, u)
		require.NoError(t, err)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		require.NoError(t, err)

		//		sampleRes := `
		//		[
		//   {
		//      "breeds":[
		//         {
		//            "weight":{
		//               "imperial":"7 - 14",
		//               "metric":"3 - 6"
		//            },
		//            "id":"esho",
		//            "name":"Exotic Shorthair",
		//            "cfa_url":"http://cfa.org/Breeds/BreedsCJ/Exotic.aspx",
		//            "vetstreet_url":"http://www.vetstreet.com/cats/exotic-shorthair",
		//            "vcahospitals_url":"https://vcahospitals.com/know-your-pet/cat-breeds/exotic-shorthair",
		//            "temperament":"Affectionate, Sweet, Loyal, Quiet, Peaceful",
		//            "origin":"United States",
		//            "country_codes":"US",
		//            "country_code":"US",
		//            "description":"The Exotic Shorthair is a gentle friendly cat that has the same personality as the Persian. They love having fun, don’t mind the company of other cats and dogs, also love to curl up for a sleep in a safe place. Exotics love their own people, but around strangers they are cautious at first. Given time, they usually warm up to visitors.",
		//            "life_span":"12 - 15",
		//            "indoor":0,
		//            "lap":1,
		//            "alt_names":"Exotic",
		//            "adaptability":5,
		//            "affection_level":5,
		//            "child_friendly":3,
		//            "dog_friendly":3,
		//            "energy_level":3,
		//            "grooming":2,
		//            "health_issues":3,
		//            "intelligence":3,
		//            "shedding_level":2,
		//            "social_needs":4,
		//            "stranger_friendly":2,
		//            "vocalisation":1,
		//            "experimental":0,
		//            "hairless":0,
		//            "natural":0,
		//            "rare":0,
		//            "rex":0,
		//            "suppressed_tail":0,
		//            "short_legs":0,
		//            "wikipedia_url":"https://en.wikipedia.org/wiki/Exotic_Shorthair",
		//            "hypoallergenic":0
		//         }
		//      ],
		//      "id":"y61B6bFCh",
		//      "url":"https://cdn2.thecatapi.com/images/y61B6bFCh.jpg",
		//      "width":898,
		//      "height":900
		//   }
		//]
		//		`

		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("some-invalid-response")),
		}
		d.EXPECT().Do(r).Return(&res, nil)
		i, err := s.GetImage(ctx)

		assert.Empty(t, i)
		assert.Error(t, err)
	})
	t.Run("returns an error as it returns a valid response with no length", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := "someurl"
		ctx := context.Background()
		d := mockcat.NewMockDoer(ctrl)
		s, err := cat.NewImageService(d, u)
		require.NoError(t, err)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		require.NoError(t, err)

		//		sampleRes := `
		//		[
		//   {
		//      "breeds":[
		//         {
		//            "weight":{
		//               "imperial":"7 - 14",
		//               "metric":"3 - 6"
		//            },
		//            "id":"esho",
		//            "name":"Exotic Shorthair",
		//            "cfa_url":"http://cfa.org/Breeds/BreedsCJ/Exotic.aspx",
		//            "vetstreet_url":"http://www.vetstreet.com/cats/exotic-shorthair",
		//            "vcahospitals_url":"https://vcahospitals.com/know-your-pet/cat-breeds/exotic-shorthair",
		//            "temperament":"Affectionate, Sweet, Loyal, Quiet, Peaceful",
		//            "origin":"United States",
		//            "country_codes":"US",
		//            "country_code":"US",
		//            "description":"The Exotic Shorthair is a gentle friendly cat that has the same personality as the Persian. They love having fun, don’t mind the company of other cats and dogs, also love to curl up for a sleep in a safe place. Exotics love their own people, but around strangers they are cautious at first. Given time, they usually warm up to visitors.",
		//            "life_span":"12 - 15",
		//            "indoor":0,
		//            "lap":1,
		//            "alt_names":"Exotic",
		//            "adaptability":5,
		//            "affection_level":5,
		//            "child_friendly":3,
		//            "dog_friendly":3,
		//            "energy_level":3,
		//            "grooming":2,
		//            "health_issues":3,
		//            "intelligence":3,
		//            "shedding_level":2,
		//            "social_needs":4,
		//            "stranger_friendly":2,
		//            "vocalisation":1,
		//            "experimental":0,
		//            "hairless":0,
		//            "natural":0,
		//            "rare":0,
		//            "rex":0,
		//            "suppressed_tail":0,
		//            "short_legs":0,
		//            "wikipedia_url":"https://en.wikipedia.org/wiki/Exotic_Shorthair",
		//            "hypoallergenic":0
		//         }
		//      ],
		//      "id":"y61B6bFCh",
		//      "url":"https://cdn2.thecatapi.com/images/y61B6bFCh.jpg",
		//      "width":898,
		//      "height":900
		//   }
		//]
		//		`

		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("[]")),
		}
		d.EXPECT().Do(r).Return(&res, nil)
		i, err := s.GetImage(ctx)

		assert.Empty(t, i)
		assert.Error(t, err)
	})

	t.Run("returns an error as it returns a valid response with no length", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := "someurl"
		ctx := context.Background()
		d := mockcat.NewMockDoer(ctrl)
		s, err := cat.NewImageService(d, u)
		require.NoError(t, err)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		require.NoError(t, err)

		sampleRes := `
				[
		  {
		     "breeds":[
		        {
		           "weight":{
		              "imperial":"7 - 14",
		              "metric":"3 - 6"
		           },
		           "id":"esho",
		           "name":"Exotic Shorthair",
		           "cfa_url":"http://cfa.org/Breeds/BreedsCJ/Exotic.aspx",
		           "vetstreet_url":"http://www.vetstreet.com/cats/exotic-shorthair",
		           "vcahospitals_url":"https://vcahospitals.com/know-your-pet/cat-breeds/exotic-shorthair",
		           "temperament":"Affectionate, Sweet, Loyal, Quiet, Peaceful",
		           "origin":"United States",
		           "country_codes":"US",
		           "country_code":"US",
		           "description":"The Exotic Shorthair is a gentle friendly cat that has the same personality as the Persian. They love having fun, don’t mind the company of other cats and dogs, also love to curl up for a sleep in a safe place. Exotics love their own people, but around strangers they are cautious at first. Given time, they usually warm up to visitors.",
		           "life_span":"12 - 15",
		           "indoor":0,
		           "lap":1,
		           "alt_names":"Exotic",
		           "adaptability":5,
		           "affection_level":5,
		           "child_friendly":3,
		           "dog_friendly":3,
		           "energy_level":3,
		           "grooming":2,
		           "health_issues":3,
		           "intelligence":3,
		           "shedding_level":2,
		           "social_needs":4,
		           "stranger_friendly":2,
		           "vocalisation":1,
		           "experimental":0,
		           "hairless":0,
		           "natural":0,
		           "rare":0,
		           "rex":0,
		           "suppressed_tail":0,
		           "short_legs":0,
		           "wikipedia_url":"https://en.wikipedia.org/wiki/Exotic_Shorthair",
		           "hypoallergenic":0
		        }
		     ],
		     "id":"y61B6bFCh",
		     "url":"https://cdn2.thecatapi.com/images/y61B6bFCh.jpg",
		     "width":898,
		     "height":900
		  }
		]
				`

		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(sampleRes)),
		}
		d.EXPECT().Do(r).Return(&res, nil)
		i, err := s.GetImage(ctx)

		assert.NoError(t, err)
		assert.Equal(t, cat.ImageURL("https://cdn2.thecatapi.com/images/y61B6bFCh.jpg"), i)
	})
}
