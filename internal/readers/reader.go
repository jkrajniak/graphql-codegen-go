package readers

import (
	"github.com/pkg/errors"
	"strings"
)

type SchemaReader interface {
	Read() ([]byte, error)
}

func DiscoverReader(schemaPath string) SchemaReader {
	if strings.Contains(schemaPath, ":") {
		return NewGitReader(schemaPath)
	} else {
		return NewLocalReader(schemaPath)
	}
}


func ReadSchemas(schemaPaths []string) ([]byte, error) {
	var outs []byte
	for _, s := range schemaPaths {
		r := DiscoverReader(s)
		o, err := r.Read()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read from %s", s)
		}
		outs = append(outs, o...)
	}
	return outs, nil
}