package readers

import (
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"regexp"
)

var repoExp, _ = regexp.Compile(`(?P<repoURL>.*\.git)(?P<path>[^@]+)@?(?P<branch>.*)`)

type GitReader struct {
	path       string
	repoURL    string
	branchName string
}

func NewGitReader(path string) *GitReader {
	params := getParams(repoExp, path)
	branchName := "master"
	if params["branch"] != "" {
		branchName = params["branch"]
	}
	return &GitReader{
		path:       params["path"],
		repoURL:    params["repoURL"],
		branchName: branchName,
	}
}

func (g *GitReader) Read() ([]byte, error) {
	memFS := memfs.New()
	r, err := git.Clone(memory.NewStorage(), memFS, &git.CloneOptions{
		URL: g.repoURL,
	})
	if err != nil {
		return nil, err
	}

	w, _ := r.Worktree()
	if err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(g.branchName)}); err != nil {
		return nil, err
	}

	file, err := memFS.Open(g.path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

func getParams(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}
