package b

import (
	"fmt"

	"github.com/qawatake/derrors"
	de "github.com/qawatake/derrors"
)

func Good() (err error) { // ok because Good begins with derrors.Wrap.
	defer derrors.Wrap(&err, "x")
	return nil
}

func Bad_simple() error { // want "should call"
	return nil
}

func private_func_is_ok() error { // ok because f is private.
	return nil
}

type S struct{}

func (s *S) Good_method() (err error) { // ok because Good_method begins with derrors.Wrap.
	defer derrors.Wrap(&err, "x")
	return nil
}

func (s *S) Bad_method() error { // want "should call"
	return nil
}

func Good_alias() (err error) { // ok because Good_alias begins with defer de.Wrap where de is an alias of derrors.
	defer de.Wrap(&err, "x")
	return nil
}

func Bad_defer_another() error { // want "should call"
	defer fmt.Println()
	return nil
}

func Good_multi_returns() (_ int, err error) { // ok because Good_multi_returns begins with defer derrors.Wrap.
	defer derrors.Wrap(&err, "x")
	return 0, nil
}

func Return_no_error() int { // ok because Return_no_error does not return error.
	return 0
}

func EmptyBody() error // ok because EmptyBody does not have body.

//lint:ignore dwrap reason
func Ignored() error { // ok because Ignored is ignored by dwrap.
	return nil
}

// lint:ignore anotherlinter reason
func Bad_another_ignore() error { // want "should call"
	return nil
}

func Bad_not_defer() (err error) { // want "should call"
	derrors.Wrap(&err, "x")
	return nil
}

func Good_not_defer_no_error() { // ok because f does not return error.
	derrors.Wrap(nil, "x")
}

func Bad_same_name() (err error) { // want "should call"
	defer Wrap(&err, "x")
	return nil
}

func Wrap(errp *error, msg string) {}
