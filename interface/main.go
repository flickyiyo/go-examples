package main

import "fmt"

func a(sd interface{}) interface{} {
	return "12"
}

func main() {
	fmt.Println(a("ads"))
}
