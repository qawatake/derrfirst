package derrors

import "fmt"

func Wrap(errp *error, msg string) {
	*errp = fmt.Errorf("%s: %w", msg, *errp)
}
