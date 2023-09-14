package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/Users/iron/.kube/config")
	if err != nil {
		log.Printf("build config from flags: %v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Printf("new for config: %v", err)
		return
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list: %v", err)
		return
	}
	fmt.Printf("there are %d pods in the clusters\n", len(pods.Items))

	for i, pod := range pods.Items {
		fmt.Printf("%d: %s, %s, %s, %s\n", i, pod.ObjectMeta.Name, pod.Status.Phase, pod.Status.HostIP, pod.Status.PodIP)
	}
}
