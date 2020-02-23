package internal

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

var s = `
overwrite: true
schema:
    - ./schema.graphql

generates:
    internal/models.go:
        config:
            packageName: internal
            entities:
                - User
                - Person
    internal/abc/models.go:
        config:
            packageName: abc
            entities:
                - Action
`

func TestReadSchema(t *testing.T) {
	stringReader := strings.NewReader(s)
	c, err := ReadConfigFromFile(stringReader)
	assert.NoError(t, err)
	sort.Slice(c.Outputs, func(i, j int) bool {
		return c.Outputs[i].PackageName < c.Outputs[j].PackageName
	})
	assert.Equal(t, []string{"./schema.graphql"}, c.Schemas)
	assert.Equal(t, "internal", c.Outputs[1].PackageName)
	assert.Equal(t, []string{"User", "Person"}, c.Outputs[1].Entities)
	assert.Equal(t, "abc", c.Outputs[0].PackageName)
	assert.Equal(t, []string{"Action"}, c.Outputs[0].Entities)
}
