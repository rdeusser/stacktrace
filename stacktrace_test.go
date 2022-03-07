package stacktrace

import (
	"errors"
	"testing"

	"github.com/hexops/autogold"
)

func caller() error {
	err := errors.New("json: cannot unmarshal json array to something or other")
	return Propagate(err, "unmarshaling json")
}

func callee() error {
	err := caller()
	return Propagate(err, "doing something important")
}

func TestPropagate(t *testing.T) {
	err := callee()
	autogold.Equal(t, autogold.Raw(err.Error()))
}
