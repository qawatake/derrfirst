package a

import "a/fmt"

func A1() error { // want "should call"
	defer fmt.Println()
	return nil
}
