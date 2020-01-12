package main

import (
	"flag"
	"github.com/jkrajniak/graphql-codegen-go/internal"
	"github.com/jkrajniak/graphql-codegen-go/internal/readers"
	"os"
	"path"
	"strings"
)

func main() {
	schemaFile := flag.String("schema", "", "schema file")
	entitiesString := flag.String("entities", "", "comma separated list of entities (optional)")
	outFile := flag.String("out", "", "output name")
	packageName := flag.String("packageName", "", "package name")
	flag.Parse()

	if packageName == nil || *packageName == "" {
		p, err := resolvePackageName()
		if err != nil {
			panic(err)
		}
		packageName = &p
	}

	schemaReader := readers.DiscoverReader(*schemaFile)
	of, err := schemaReader.Read()
	if err != nil {
		panic(err)
	}

	var output internal.Outputer = &internal.STDOutput{}
	if outFile != nil && *outFile != "" {
		o, err := internal.NewFileOutput(*outFile)
		if err != nil {
			panic(err)
		}
		output = o

	}

	var entities []string
	if *entitiesString != "" {
		entities = strings.Split(*entitiesString, ",")
	}

	goGenerator := internal.NewGoGenerator(output, entities, *packageName)
	goGenerator.Generate(string(of), *schemaFile)

	if err := output.Close(); err != nil {
		panic(err)
	}
}

func resolvePackageName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, packageName := path.Split(cwd)

	return packageName, nil
}
