package a

import (
	"fmt"
	fmt1 "fmt"
)

func f() { // ok
	fmt.Println("f")
}

func F() { // want "should call"
	fmt.Println("F")
}

func g() { // ok
	defer f()
}

func G() { // want "should call"
	defer f()
}

func h() { // ok
	defer fmt.Println("a")
}

func H() { // ok
	defer fmt.Println("a")
}

func I() { // ok
	defer fmt.Println()
}

func J() { // ok
	defer fmt1.Println()
}
