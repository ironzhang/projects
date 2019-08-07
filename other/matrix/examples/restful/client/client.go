package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/httputils"
	"github.com/ironzhang/matrix/restful"
	"github.com/ironzhang/matrix/tlog"
)

const url = "http://localhost:8080"

var client = restful.Client{
	Client:  &http.Client{Transport: httputils.NewVerboseRoundTripper(nil, nil, nil)},
	Context: context_value.WithVerbose(context.Background(), true),
}

func main() {
	tlog.Init(tlog.Config{DisableStacktrace: true})

	var cmd string
	flag.StringVar(&cmd, "cmd", "root", "cmd name")
	flag.Parse()

	switch cmd {
	case "root":
		DoRootCommand()
	case "echo":
		DoEchoCommand()
	default:
		fmt.Println("unknown command.")
	}
}

func DoRootCommand() {
	log := tlog.Std().Sugar()
	var resp string
	if err := client.Get(url, nil, &resp); err != nil {
		log.Errorw("client get", "error", err)
		return
	}
	fmt.Println(resp)
}

func DoEchoCommand() {
	log := tlog.Std().Sugar()
	var req = strings.Join(flag.Args(), " ")
	var resp string
	if err := client.Post(url+"/echo", req, &resp); err != nil {
		log.Errorw("client post", "error", err)
		return
	}
	fmt.Println(resp)
}
