package main

import (
	"github.com/ironzhang/tlog"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	tlog.Info(dmp.DiffPrettyText(diffs))
}
