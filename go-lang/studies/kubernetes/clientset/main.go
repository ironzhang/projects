package main

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Printf("build config from flags: %v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Printf("new client set: %v", err)
		return
	}

	slist, err := clientset.CoreV1().Services("dev").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list services: %v", err)
		return
	}
	for _, svc := range slist.Items {
		selector := labels.NewSelector()
		for key, value := range svc.Spec.Selector {
			r, err := labels.NewRequirement(key, selection.Equals, []string{value})
			if err != nil {
				log.Printf("new requirement: %v", err)
				return
			}
			selector.Add(*r)
		}

		plist, err := clientset.CoreV1().Pods("dev").List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			log.Printf("list pods: %v", err)
			return
		}
		for _, pod := range plist.Items {
			for _, c := range pod.Spec.Containers {
				for _, port := range c.Ports {
					log.Printf("%s: %s:%d", port.Name, pod.Status.PodIP, port.ContainerPort)
				}
			}
		}
	}
}
