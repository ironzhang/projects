package main

import (
	"fmt"
	"log"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Printf("build config from flags: %v", err)
		return
	}

	discoveryclient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		log.Printf("new discovery client: %v", err)
		return
	}

	lists, err := discoveryclient.ServerResources()
	if err != nil {
		log.Printf("get server resources: %v", err)
		return
	}
	for _, list := range lists {
		for _, apires := range list.APIResources {
			fmt.Printf("%s,%s\n", apires.Name, apires.Kind)
		}
	}
}
