package a

import "fmt"

func f() { // want "should call"
	fmt.Println("f")
}

func g() { // want "should call"
	defer f()
}

func h() { // ok
	defer fmt.Println("a")
}

func i() { // ok
	defer fmt.Println()
}
