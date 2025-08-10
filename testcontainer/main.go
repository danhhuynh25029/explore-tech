package main

import "fmt"

func main() {
	type B struct {
		B string
	}
	type A struct {
		A string
		B *B
	}

	var a *A

	a = &A{
		A: "a",
		B: &B{B: "b"},
	}
	fmt.Println(fmt.Sprintf("%+v", a.B))
}
