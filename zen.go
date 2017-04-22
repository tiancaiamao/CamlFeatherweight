package main

import (
	"fmt"

	"github.com/tiancaiamao/CamlFeatherweight/zen"
)

func main() {
	vm, err := zen.LoadFile("test.out")
	if err != nil {
		panic(err)
	}
	v := vm.Run()
	fmt.Println(v)
}
