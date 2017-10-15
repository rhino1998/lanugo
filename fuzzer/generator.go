package fuzzer

import (
	"io"
	"strings"

	"github.com/google/skylark"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

//Generator wraps a skylark.Callable and puts its output into an io.ReadCloser
type Generator struct {
	fn skylark.Callable
}

//NewGenerator returns a new Generator from a skylark.Callable
func NewGenerator(fn skylark.Callable) *Generator {
	return &Generator{fn: fn}
}

//Generate returns an io.ReadCloser based on the skylark.Callable value it wraps
func (g *Generator) Generate() (io.ReadCloser, error) {
	thread := &skylark.Thread{}
	val, err := g.fn.Call(thread, nil, nil)

	if err != nil {
		return nil, err
	}

	switch vt := val.(type) {
	case *ReadCloserValue:
		return vt, nil
	case skylark.String:
		return nopCloser{strings.NewReader(string(vt))}, nil
	default:
		return nopCloser{strings.NewReader(vt.String())}, nil
	}
}