package main

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func PrintPod(pod *v1.Pod) {
	for _, c := range pod.Spec.Containers {
		for _, cp := range c.Ports {
			if pod.Status.PodIP != "" {
				fmt.Printf("%s/%s\t%s.%s\t%s\t%s:%d\n",
					pod.ObjectMeta.Namespace, pod.ObjectMeta.Name,
					cp.Name, pod.ObjectMeta.Labels["app"],
					pod.Status.Phase, pod.Status.PodIP, cp.ContainerPort)
			}
		}
	}
}

type PodEventHandler struct {
	queue workqueue.RateLimitingInterface
}

func (p *PodEventHandler) OnAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		p.queue.Add(key)
	}
}

func (p *PodEventHandler) OnUpdate(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err == nil {
		p.queue.Add(key)
	}
}

func (p *PodEventHandler) OnDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		p.queue.Add(key)
	}
}

type Controller struct {
	queue    workqueue.RateLimitingInterface
	indexer  cache.Indexer
	informer cache.Controller
}

func (p *Controller) process() bool {
	key, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(key)

	log.Printf("process %s", key)
	objs := p.indexer.List()
	for _, obj := range objs {
		pod := obj.(*v1.Pod)
		PrintPod(pod)
	}

	return true
}

func (p *Controller) Run(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer p.queue.ShutDown()

	go p.informer.Run(stopCh)
	if !cache.WaitForNamedCacheSync("controller", stopCh, p.informer.HasSynced) {
		log.Printf("timed out waiting for caches to sync")
		return
	}
	for p.process() {
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

	// new workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	lselector := labels.NewSelector()
	r, err := labels.NewRequirement("app", selection.Equals, []string{"nginx"})
	if err != nil {
		log.Printf("new requirement: %v", err)
		return
	}
	lselector = lselector.Add(*r)
	fselector := fields.Everything()

	podEventHandler := &PodEventHandler{queue: queue}
	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = lselector.String()
		options.FieldSelector = fselector.String()
	}
	podsListWatcher := cache.NewFilteredListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", "dev", optionsModifier)
	indexer, informer := cache.NewIndexerInformer(podsListWatcher, &v1.Pod{}, 0, podEventHandler, cache.Indexers{})

	ch := make(chan struct{})
	c := Controller{queue: queue, indexer: indexer, informer: informer}
	c.Run(ch)
}
