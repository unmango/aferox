package gitignore

import (
	"bufio"
	"fmt"
	"io"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
	"github.com/unmango/aferox/filter"
)

const DefaultFile = ".gitignore"

type Ignore = ignore.IgnoreParser

func NewFsFromLines(base afero.Fs, lines ...string) afero.Fs {
	return NewFsFromIgnore(base, ignore.CompileIgnoreLines(lines...))
}

func NewFsFromIgnore(base afero.Fs, ignore Ignore) afero.Fs {
	return filter.NewFs(base, func(s string) bool {
		return !ignore.MatchesPath(s)
	})
}

func NewFsFromReader(base afero.Fs, reader io.Reader) (afero.Fs, error) {
	lines := []string{}
	s := bufio.NewScanner(reader)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("scanning ignore lines: %w", s.Err())
	}

	return NewFsFromLines(base, lines...), nil
}

func NewFsFromFile(base afero.Fs, path string) (afero.Fs, error) {
	f, err := base.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening ignore file: %w", err)
	}
	defer f.Close()
	return NewFsFromReader(base, f)
}

func OpenDefault(base afero.Fs) (afero.Fs, error) {
	return NewFsFromFile(base, DefaultFile)
}
