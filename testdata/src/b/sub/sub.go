package sub

import "github.com/qawatake/derrors"

func Good() (err error) { // ok because Good begins with derrors.Wrap.
	defer derrors.Wrap(&err, "x")
	return nil
}

func Bad_simple() error { // want "should call"
	return nil
}
