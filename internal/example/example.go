package example

import (
	"errors"
	"fmt"
)

func Good() (err error) {
	defer Wrap(&err)
	err = doSomething()
	return nil
}

func Bad() error { // <-  should call defer Wrap fist.
	return doSomething()
}

//lint:ignore dwrap this is because ...
func Ignored() error {
	return doSomething()
}

func doSomething() error {
	return errors.New("original error")
}

func Wrap(errp *error) {
	if *errp != nil {
		*errp = fmt.Errorf("wrapped: %w", *errp)
	}
}
