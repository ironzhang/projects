package main

import (
	"github.com/ironzhang/matrix/framework"
	_ "github.com/ironzhang/matrix/framework/modules/dashboard-module"
	_ "github.com/ironzhang/matrix/framework/modules/debug-module"
)

func main() {
	framework.Main()
}
