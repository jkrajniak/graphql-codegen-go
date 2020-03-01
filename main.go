package main

import (
	"flag"
	"github.com/jkrajniak/graphql-codegen-go/internal"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func main() {
	configYaml := flag.String("config", "", "config yaml")
	schemasFile := flag.String("schemas", "", "schema file (comma separated list)")
	entitiesString := flag.String("entities", "", "comma separated list of entities (optional)")
	packageNameString := flag.String("packageName", "", "package name")
	outFile := flag.String("out", "", "file output name (optional, default: stdout)")
	flag.Parse()

	var config internal.Config
	if configYaml != nil && *configYaml != "" {
		configFile, err := os.Open(*configYaml)
		if err != nil {
			panic(err)
		}
		defer configFile.Close()

		yc, err := internal.ReadConfigFromFile(configFile)
		if err != nil {
			panic(err)
		}
		config = yc
	} else {
		goGenerateDate := internal.GetGOGenerate()

		var entities []string
		if *entitiesString != "" {
			entities = strings.Split(*entitiesString, ",")
		}

		var pkgName *string
		if goGenerateDate != nil {
			pkgName = &goGenerateDate.GOPackage
		}
		if packageNameString != nil && *packageNameString != "" {
			pkgName = packageNameString
		}
		if pkgName == nil {
			flag.Usage()
			panic("packageName not defined")
		}

		config = internal.Config{
			Schemas: strings.Split(*schemasFile, ","),
			Outputs: []internal.OutputItem{
				{OutputPath: *outFile, PackageName: *pkgName, Entities: entities},
			},
		}
	}

	// Combine all schemas.
	inputSchemas, err := internal.ReadSchemas(config.Schemas)
	if err != nil {
		panic(err)
	}

	for _, o := range config.Outputs {
		output, err := internal.NewFileOutput(o.OutputPath)
		if err != nil {
			panic(errors.Wrapf(err, "failed to create output to %s", o.OutputPath))
		}

		loadedDocs, err := internal.LoadSchemas(inputSchemas)
		if err != nil {
			panic(errors.Wrapf(err, "failed to parse input schemas"))
		}

		gen := internal.NewGoGenerator(output, o.Entities, o.PackageName)
		if err := gen.Generate(loadedDocs); err != nil {
			panic(errors.Wrapf(err, "failed to generate go structs"))
		}

		if err := output.Close(); err != nil {
			panic(err)
		}
	}
}
