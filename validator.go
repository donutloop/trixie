package trixie

//Validator validates the incomming value against a valid value/s
type Validator interface {
	Validate(RouteInterface) error
}

//pathValidator check if a path is set and validates the value .
type pathValidator struct{}

func NewPathValidator() pathValidator {
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

//methodValidator check if method is a correct value.
type methodValidator struct{}

func NewMethodValidator() methodValidator {
	return methodValidator{}
}

func (v methodValidator) Validate(r RouteInterface) error {

	for k := range r.GetHandlers() {
		if method := Methods.lookupID(k); method == "" {
			return NewBadMethodError()
		}
	}

	return nil
}

var Validatoren = []Validator{
	NewPathValidator(),
	NewMethodValidator(),
}
