package tmux

//Validator validates the incomming value against a valid value/s
type Validator interface {
	Validate(RouteInterface) error
}

//pathValidator check if a path is set and validates the value .
type pathValidator struct{}

func newPathValidator() pathValidator {
	return pathValidator{}
}

func (v pathValidator) Validate(r RouteInterface) error {

	if len(r.GetPattern()) == 0 {
		return NewBadPathError("Path is empty")
	}

	if r.GetPattern()[0] != '/' {
		return NewBadPathError("Path starts not with a /")
	}

	return nil
}

var Validatoren = []Validator{
	newPathValidator(),
}
