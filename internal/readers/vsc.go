package readers

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"regexp"
)

var repoExp, _ = regexp.Compile(`(?P<repoURL>.*\.git)(?P<path>[^(@|#)?]+)(?P<hashBranch>(@|#)?)(?P<ref>.*)`)

type GitReader struct {
	path            string
	repoURL         string
	refName         string
	branchName      string
	isRefCommitHash bool
}

func NewGitReader(path string) *GitReader {
	params := getParams(repoExp, path)
	branchName := "master"
	refName := ""
	isRefCommitHash := false
	if params["hashBranch"] == "#" {
		refName = params["ref"]
		isRefCommitHash = true
	} else if params["hashBranch"] == "@" {
		branchName = params["ref"]
	}

	return &GitReader{
		path:            params["path"],
		repoURL:         params["repoURL"],
		refName:         refName,
		branchName:      branchName,
		isRefCommitHash: isRefCommitHash,
	}
}

func (g *GitReader) Read() ([]byte, error) {
	memFS := memfs.New()
	var r *git.Repository
	var err error

	r, err = git.Clone(memory.NewStorage(), memFS, &git.CloneOptions{
		URL:           g.repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(g.branchName),
	})
	if err != nil {
		if err.Error() == "reference not found" {
			r, err = git.Clone(memory.NewStorage(), memFS, &git.CloneOptions{
				URL:           g.repoURL,
				ReferenceName: plumbing.NewTagReferenceName(g.branchName),
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Wrapf(err, "clone with branch %s", g.branchName)
		}
	}

	if g.isRefCommitHash {
		w, _ := r.Worktree()
		if err = w.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(g.refName)}); err != nil {
			return nil, errors.Wrapf(err, "checkout with ref %s", g.refName)
		}
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
