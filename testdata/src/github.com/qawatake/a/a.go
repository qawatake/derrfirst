package a

func Good() (err error) { // ok because Good begins with derrors.Wrap.
	defer Wrap(&err, "x")
	return nil
}

func Bad_simple() error { // want "should call"
	return nil
}

func Wrap(errp *error, msg string) {}
