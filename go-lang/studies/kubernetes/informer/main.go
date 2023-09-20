package main

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintPods(pods []*v1.Pod) {
	for _, pod := range pods {
		fmt.Printf("%s\n", pod.String())
	}
}

type Controller struct {
	lister   corev1.PodLister
	informer cache.SharedIndexInformer
}

func NewController(lister corev1.PodLister, informer cache.SharedIndexInformer) *Controller {
	return &Controller{
		lister:   lister,
		informer: informer,
	}
}

func (p *Controller) Run() error {
	selector := labels.NewSelector()
	requirement, err := labels.NewRequirement("app", selection.Equals, []string{"nginx"})
	if err != nil {
		return err
	}
	selector.Add(*requirement)

	pods, err := p.lister.Pods("dev").List(selector)
	if err != nil {
		return err
	}
	PrintPods(pods)

	return nil
}

func main() {
	// build config
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Printf("build config from flags: %v", err)
		return
	}

	// new client
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Printf("new client set: %v", err)
		return
	}

	// new informer
	factory := informers.NewSharedInformerFactory(clientset, 0)
	podInformer := factory.Core().V1().Pods()

	// run controller
	controller := NewController(podInformer.Lister(), podInformer.Informer())
	if err = controller.Run(); err != nil {
		log.Printf("controller run: %v", err)
		return
	}
}
