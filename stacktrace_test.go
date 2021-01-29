package stacktrace

import (
	"errors"
	"testing"

	"github.com/hexops/autogold"
)

func fnlibcalls() error {
	err := errors.New("json: cannot unmarshal json array to something or other")
	return Propagate(err, "unmarshaling json")
}

func fninlib() error {
	err := fnlibcalls()
	return Propagate(err, "doing something important")
}

func TestPropogate(t *testing.T) {
	err := fninlib()
	autogold.Equal(t, autogold.Raw(err.Error()))
}
