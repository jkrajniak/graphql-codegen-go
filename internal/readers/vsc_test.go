package readers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGitReader(t *testing.T) {
	s := NewGitReader("https://github.com/test/repo.git/internal/schema.sql")
	assert.Equal(t, "master", s.branchName)
	assert.Equal(t, "/internal/schema.sql", s.path)
	assert.Equal(t, "https://github.com/test/repo.git", s.repoURL)
}


func TestNewGitReader_BranchName(t *testing.T) {
	s := NewGitReader("https://github.com/test/repo.git/internal/schema.sql@branch")
	assert.Equal(t, "branch", s.branchName)
	assert.Equal(t, "/internal/schema.sql", s.path)
	assert.Equal(t, "https://github.com/test/repo.git", s.repoURL)
}
