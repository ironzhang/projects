package main

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintPod(pod *v1.Pod) {
	fmt.Printf("%s.%s\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
}

type PodEventHandler struct {
	store cache.Store
}

func (p *PodEventHandler) OnAdd(obj interface{}) {
	objs := p.store.List()
	if len(objs) > 0 {
		log.Printf("on add:")
		for _, obj := range objs {
			pod := obj.(*v1.Pod)
			PrintPod(pod)
		}
	}
}

func (p *PodEventHandler) OnUpdate(oldObj, newObj interface{}) {
	objs := p.store.List()
	if len(objs) > 0 {
		log.Printf("on update:")
		for _, obj := range objs {
			pod := obj.(*v1.Pod)
			PrintPod(pod)
		}
	}
}

func (p *PodEventHandler) OnDelete(obj interface{}) {
	objs := p.store.List()
	if len(objs) > 0 {
		log.Printf("on delete:")
		for _, obj := range objs {
			pod := obj.(*v1.Pod)
			PrintPod(pod)
		}
	}
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

	lselector := labels.NewSelector()
	r, err := labels.NewRequirement("app", selection.Equals, []string{"nginx"})
	if err != nil {
		log.Printf("new requirement: %v", err)
		return
	}
	lselector = lselector.Add(*r)
	fselector := fields.Everything()

	podEventHandler := &PodEventHandler{}
	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = lselector.String()
		options.FieldSelector = fselector.String()
	}
	podsListWatcher := cache.NewFilteredListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", "dev", optionsModifier)
	store, controller := cache.NewInformer(podsListWatcher, &v1.Pod{}, 0, podEventHandler)
	podEventHandler.store = store

	ch := make(chan struct{})
	controller.Run(ch)
}
