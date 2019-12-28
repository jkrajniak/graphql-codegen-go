package main

import (
	"flag"
	"github.com/jkrajniak/graphql-codegen-go/internal"
	"io/ioutil"
	"strings"
)

func main() {
	schemaFile := flag.String("schema", "", "schema file")
	entitiesString := flag.String("entities", "", "comma separated list of entities (optional)")
	flag.Parse()

	of, err := ioutil.ReadFile(*schemaFile)
	if err != nil {
		panic(err)
	}

	output := internal.STDOutput{}

	var entities []string
	if *entitiesString != "" {
		entities = strings.Split(*entitiesString, ",")
	}

	goGenerator := internal.NewGoGenerator(&output, entities)
	goGenerator.Generate(string(of))
}
