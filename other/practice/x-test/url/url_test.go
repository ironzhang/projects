package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	var err error
	var path1 = "https://open.bot.tmall.com/oauth/callback"
	var path2 = "https%3A%2F%2Fopen.bot.tmall.com%2Foauth%2Fcallback"

	path := url.QueryEscape(path1)
	if path != path2 {
		t.Fatalf("%s != %s", path, path2)
	}
	fmt.Println(path)

	path, err = url.QueryUnescape(path1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)

	path, err = url.QueryUnescape(path2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
}
