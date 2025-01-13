package gfs

import (
	"errors"
	"fmt"
	"io/fs"
	"reflect"

	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

type containFileWithBytes struct {
	path  string
	bytes []byte
}

// Match implements types.GomegaMatcher.
func (c *containFileWithBytes) Match(actual interface{}) (success bool, err error) {
	fs, ok := actual.(afero.Fs)
	if !ok {
		return false, fmt.Errorf("expected an [afero.Fs] got %s", reflect.TypeOf(actual))
	}

	return afero.FileContainsBytes(fs, c.path, c.bytes)
}

// FailureMessage implements types.GomegaMatcher.
func (c *containFileWithBytes) FailureMessage(actual interface{}) (message string) {
	fs, ok := actual.(afero.Fs)
	if !ok {
		return fmt.Sprintf("expected an [afero.Fs] got %s", reflect.TypeOf(actual))
	}

	data, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf(
		"expected file at\n%s\n\tto contain content:\n%s\n\tbut instead had\n%s",
		c.path, c.bytes, data,
	)
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (c *containFileWithBytes) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected file at\n%s\n\tnot to contain content:\n%s", c.path, c.bytes)
}

func ContainFileWithBytes(path string, bytes []byte) types.GomegaMatcher {
	return &containFileWithBytes{path, bytes}
}

type containFile struct {
	path string
}

// Match implements types.GomegaMatcher.
func (c *containFile) Match(actual interface{}) (success bool, err error) {
	fs, ok := actual.(afero.Fs)
	if !ok {
		return false, fmt.Errorf("expected an [afero.Fs] got %s", reflect.TypeOf(actual))
	}

	_, err = fs.Open(c.path)
	return err == nil, nil
}

// FailureMessage implements types.GomegaMatcher.
func (c *containFile) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected file to exist at %s", c.path)
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (c *containFile) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected %s not to exist", c.path)
}

func ContainFile(path string) types.GomegaMatcher {
	return &containFile{path}
}

type beEquivalentToFs struct {
	expected afero.Fs
}

// Match implements types.GomegaMatcher.
func (e *beEquivalentToFs) Match(actual interface{}) (success bool, err error) {
	target, ok := actual.(afero.Fs)
	if !ok {
		return false, fmt.Errorf("exected an [afero.Fs] but got %s", reflect.TypeOf(actual))
	}

	failures := []error{}
	err = afero.Walk(e.expected, "",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				exists, err := afero.DirExists(target, path)
				if err != nil {
					return err
				}
				if !exists {
					failures = append(failures,
						fmt.Errorf("expected dir to exist at %s", path),
					)
				}

				return nil
			}

			exists, err := afero.Exists(target, path)
			if err != nil {
				return err
			}
			if !exists {
				failures = append(failures,
					fmt.Errorf("expected file to exist at %s", path),
				)

				return nil
			}

			expectedBytes, err := afero.ReadFile(e.expected, path)
			if err != nil {
				return err
			}

			matched, err := afero.FileContainsBytes(target, path, expectedBytes)
			if err != nil {
				return err
			}
			if !matched {
				actualBytes, err := afero.ReadFile(target, path)
				if err != nil {
					return err
				}

				failures = append(failures,
					fmt.Errorf("expected file at %s to contain content:\n\t%s\nbut found\n\t%s",
						path, string(expectedBytes), string(actualBytes),
					),
				)
			}

			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("walking expected filesystem: %w", err)
	}

	return len(failures) == 0, errors.Join(failures...)
}

// FailureMessage implements types.GomegaMatcher.
func (e *beEquivalentToFs) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"expected fs %s to match fs %s",
		actual.(afero.Fs).Name(),
		e.expected.Name(),
	)
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (e *beEquivalentToFs) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"expected fs %s not to match fs %s",
		actual.(afero.Fs).Name(),
		e.expected.Name(),
	)
}

func BeEquivalentToFs(fs afero.Fs) types.GomegaMatcher {
	return &beEquivalentToFs{fs}
}
