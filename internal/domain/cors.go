package domain

type Cors struct {
	Id     int
	Name   string
	Origin string
}

func NewCors(name string, origin string) (*Cors, error) {
	cors := Cors{
		Name:   name,
		Origin: origin,
	}

	return &cors, nil
}
