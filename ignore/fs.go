package ignore

import (
	"bufio"
	"fmt"
	"io"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
)

const DefaultFile = ".gitignore"

type Ignore interface {
	MatchesPath(string) bool
}

func NewFsFromGitIgnoreLines(base afero.Fs, lines ...string) afero.Fs {
	return NewFsFromIgnore(base, ignore.CompileIgnoreLines(lines...))
}

func NewFsFromIgnore(base afero.Fs, ignore Ignore) afero.Fs {
	return filter.NewFs(base, not(ignore.MatchesPath))
}

func NewFsFromGitIgnoreReader(base afero.Fs, reader io.Reader) (afero.Fs, error) {
	lines := []string{}
	s := bufio.NewScanner(reader)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("scanning ignore lines: %w", s.Err())
	}

	return NewFsFromGitIgnoreLines(base, lines...), nil
}

func NewFsFromGitIgnoreFile(base afero.Fs, path string) (afero.Fs, error) {
	if f, err := base.Open(path); err != nil {
		return nil, fmt.Errorf("opening ignore file: %w", err)
	} else {
		defer f.Close()
		return NewFsFromGitIgnoreReader(base, f)
	}
}

func OpenDefaultGitIgnore(base afero.Fs) (afero.Fs, error) {
	return NewFsFromGitIgnoreFile(base, DefaultFile)
}

// This is entirely unnecessary
func not(fn func(string) bool) func(string) bool {
	return func(s string) bool {
		return !fn(s)
	}
}
