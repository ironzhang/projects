package main

import (
	_ "net/http"
)

func init() {
	println("main.init.2")
}

func main() {
	test()
}

func init() {
	println("main.init.1")
}
