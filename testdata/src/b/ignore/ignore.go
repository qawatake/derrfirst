package ignore

func Ignored() error { // ok because ignore package is ignored.
	return nil
}
