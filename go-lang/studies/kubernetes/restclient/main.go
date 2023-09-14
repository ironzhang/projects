package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Printf("build config from flags: %v", err)
		return
	}
	cfg.GroupVersion = &v1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs
	cfg.APIPath = "/api"

	restcli, err := rest.RESTClientFor(cfg)
	if err != nil {
		log.Printf("rest client for: %v", err)
		return
	}

	var list v1.PodList
	err = restcli.Get().Resource("pods").Namespace("dev").Do(context.TODO()).Into(&list)
	if err != nil {
		log.Printf("get pods: %v", err)
		return
	}
	for _, pod := range list.Items {
		fmt.Printf("%s,%s,%s\n", pod.Name, pod.Status.HostIP, pod.Status.PodIP)
	}
}
