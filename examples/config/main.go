package main

import (
	"fmt"
	"github.com/jkrajniak/graphql-codegen-go/examples/config/internal"
)

func main() {
	person := internal.Person{
		FirstName: "",
		Lastname:  "",
		Age:       nil,
		Gender:    nil,
		Address:   nil,
	}
	fmt.Println(person)
}