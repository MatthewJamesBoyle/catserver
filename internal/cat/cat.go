package cat

type Service struct {
	img  ImageGetter
	fact FactGetter
}

func NewService(getter ImageGetter, factGetter FactGetter) (*Service, error) {
	return nil, nil
}
