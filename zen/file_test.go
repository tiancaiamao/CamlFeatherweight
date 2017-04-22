package zen

import (
	"fmt"
	"testing"

	"github.com/juju/errors"
)

func TestLoadFile(t *testing.T) {
	vm, err := LoadFile("test.out")
	if err != nil {
		t.Error(errors.ErrorStack(err))
	}
	fmt.Println(len(vm.code), vm.code)
}
