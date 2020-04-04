package gen
//go:generate mockgen -package mockcat -destination internal/mock/mockcat/cat.go github.com/matthewjamesboyle/catserver/internal/cat FactGetter,ImageGetter,Doer

