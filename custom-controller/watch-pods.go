package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type controller struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path of kubeconfig")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller := NewController(factory)
	stop := make(chan struct{})
	defer close(stop)

	err = controller.Run(stop)
	if err != nil {
		panic(err.Error())
	}
	select {}
}

func NewController(informerFactory informers.SharedInformerFactory) *controller {
	podInformer := informerFactory.Core().V1().Pods()
	c := &controller{
		informerFactory: informerFactory,
		podInformer:     podInformer,
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller-name"),
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("CREATED POD: %s/%s", pod.Namespace, pod.Name)
		},
		UpdateFunc: func(old, new interface{}) {
			o := old.(*v1.Pod)
			n := new.(*v1.Pod)
			fmt.Println("UPDATED POD: %s/%s %s", o.Namespace, o.Name, n.Status.Phase)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("DELETED POD: %s/%s", pod.Namespace, pod.Name)
		},
	})
	return c
}

func (c *controller) Run(stop chan struct{}) error {
	c.informerFactory.Start(stop)
	if !cache.WaitForCacheSync(stop, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}
