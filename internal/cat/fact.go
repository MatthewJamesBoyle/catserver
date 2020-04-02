package cat

import "context"

type Fact string

type FactGetter interface {
	GetFact(ctx context.Context) (Fact, error)
}
