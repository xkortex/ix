package main

import (
	"github.com/xkortex/ix/ix"
)

var (
	Version = "unset"
)

func main() {
	ix.Version = Version // todo: need less hackish way of setting version
	ix.Execute()
}
