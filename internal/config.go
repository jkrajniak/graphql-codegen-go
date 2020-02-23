package internal

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type YamlGenerateConfig struct {
	PackageName string   `yaml:"packageName"`
	Entities    []string `yaml:"entities"`
}

type YamlGenerateItem struct {
	Config YamlGenerateConfig `yaml:"config"`
}

type YamlConfig struct {
	Schema    []string                    `yaml:"schema"`
	Generates map[string]YamlGenerateItem `yaml:"generates"`
}

type OutputItem struct {
	OutputPath  string
	PackageName string
	Entities    []string
}

type GOGenerate struct {
	GOFile    string
	GOLine    int
	GOPackage string
}

type Config struct {
	Schemas   []string
	Outputs   []OutputItem

	InPlace     bool
	GoGenConfig *GOGenerate
}

func ReadConfigFromFile(f io.Reader) (Config, error) {
	configYaml, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	var yamlConfig YamlConfig
	if err := yaml.Unmarshal(configYaml, &yamlConfig); err != nil {
		return Config{}, err
	}

	c := Config{
		Schemas:   yamlConfig.Schema,
		Outputs:   []OutputItem{},
	}
	for outPath, item := range yamlConfig.Generates {
		c.Outputs = append(c.Outputs, OutputItem{
			OutputPath:  outPath,
			PackageName: item.Config.PackageName,
			Entities:    item.Config.Entities,
		})
	}

	return c, nil
}

func GetGOGenerate() *GOGenerate {
	if goFile, has := os.LookupEnv("GOFILE"); has {
		goLine, err := strconv.Atoi(os.Getenv("GOLINE"))
		if err != nil {
			panic(err)
		}
		return &GOGenerate{
			GOFile:    goFile,
			GOLine:    goLine,
			GOPackage: os.Getenv("GOPACKAGE"),
		}
	}

	return nil
}
