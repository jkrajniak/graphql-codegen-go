package readers

import "strings"

type SchemaReader interface {
	Read() ([]byte, error)
}

func DiscoverReader(schemaPath string) SchemaReader {
	if strings.Contains(schemaPath, "://") {
		return NewGitReader(schemaPath)
	} else {
		return NewLocalReader(schemaPath)
	}
}
