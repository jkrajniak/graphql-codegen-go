package main

import (
	"fmt"
	"github.com/jkrajniak/graphql-codegen-go/examples/local-schema/internal/appsync"
)

func main() {
	newYear := appsync.EnumYearNEW
	fmt.Println(newYear)
	e := appsync.Entity1{Y: &newYear}
	fmt.Printf("%v\n", e)
}
