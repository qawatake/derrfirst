package a

import (
	"fmt"
	fmt1 "fmt"
)

func f() { // ok because f is private.
	fmt.Println("f")
}

func F() { // want "should call"
	fmt.Println("F")
}

func G() { // want "should call"
	defer f()
}

func H() { // ok because h begins by deferring a call to fmt.Println.
	defer fmt.Println("a")
}

func J() { // ok because j begins by deferring a call to fmt1.Println where fmt1 is alias of fmt.
	defer fmt1.Println()
}
