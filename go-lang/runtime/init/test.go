package main

import "little.io/practice/init/sum"

func init() {
	println("test.init")
}

func test() {
	println(sum.Sum(1, 2, 3))
}
