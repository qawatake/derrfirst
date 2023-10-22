package b

import "b/derrors"

func Bad_another_pkg() (err error) { // want "should call"
	defer derrors.Wrap(&err, "x")
	return nil
}
