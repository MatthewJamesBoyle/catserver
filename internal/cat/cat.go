package cat

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type CatResult struct {
	ImageURL ImageURL
	Fact     Fact
}

type Service struct {
	img  ImageGetter
	fact FactGetter
}

type ErrNilParam struct {
	Parameter string
}

func (e ErrNilParam) Error() string {
	return fmt.Sprintf("%s was nil", e.Parameter)
}

type ErrServiceError struct {
	UnderLyingError error
}

func NewService(getter ImageGetter, factGetter FactGetter) (*Service, error) {
	if getter == nil {
		return nil, ErrNilParam{Parameter: "ImageGetter"}
	}
	if factGetter == nil {
		return nil, ErrNilParam{Parameter: "FactGetter"}
	}

	return &Service{
		img:  getter,
		fact: factGetter,
	}, nil
}

func (s *Service) GetImageAndFact(ctx context.Context) (CatResult, error) {

	eg, ctx := errgroup.WithContext(ctx)

	var f Fact
	var i ImageURL
	eg.Go(func() error {
		ft, err := s.fact.GetFact(ctx)
		f = ft
		if err != nil {
			return fmt.Errorf("GetImageAndFact GetFact: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		it, err := s.img.GetImage(ctx)
		i = it
		if err != nil {
			return fmt.Errorf("GetImageAndFact GetImage: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return CatResult{}, err
	}

	return CatResult{
		ImageURL: i,
		Fact:     f,
	}, nil
}
