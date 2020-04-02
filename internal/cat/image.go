package cat

import "context"

type ImageURL string

type ImageGetter interface {
	GetImage(ctx context.Context) (ImageURL, error)
}
