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

func TestNewGitReader_SSHPath(t *testing.T) {
	s := NewGitReader("git@github.com:jkrajniak/sc-priv.git/schema.sql")
	assert.Equal(t, "master", s.branchName)
	assert.Equal(t, "/schema.sql", s.path)
	assert.Equal(t, "git@github.com:jkrajniak/sc-priv.git", s.repoURL)
}

func TestNewGitReader_BranchName(t *testing.T) {
	s := NewGitReader("https://github.com/test/repo.git/internal/schema.sql@branch")
	assert.Equal(t, "branch", s.branchName)
	assert.Equal(t, "/internal/schema.sql", s.path)
	assert.Equal(t, "https://github.com/test/repo.git", s.repoURL)
}

func TestNewGitReader_Hash(t *testing.T) {
	s := NewGitReader("https://github.com/test/repo.git/internal/schema.sql#hash12344c1234")
	assert.Equal(t, "hash12344c1234", s.refName)
	assert.Equal(t, "/internal/schema.sql", s.path)
	assert.Equal(t, "https://github.com/test/repo.git", s.repoURL)
}
