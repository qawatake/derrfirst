package c

func Good_wrapping_func_defined_in_the_same_pkg() (err error) { // ok because the func begins with defer Wrap.
	defer Wrap(&err, "x")
	return nil
}

func Wrap(errp *error, msg string) {}
