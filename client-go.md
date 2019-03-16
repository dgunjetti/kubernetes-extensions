# Kubernetes Custom Controllers

Controllers are control loops that watches the shared state of cluster through API server and makes changes to move the current state to desired state. 
Ex: Deployment controller, Replicaset controller, node controller etc..

We create our own controllers for core and custom resources.

## Controller componentes

### SharedInformer
Informer is a shared data cache. Informer deliver events to listeners when the data they are interested in changes. We create reference to Informer using factory method.
```
podInformer = InformerFactory.Core().V1().Pods()
```

### Resource Event Handlers
Controllers registers interest in a specific object using EventHandlers. Controller can handle create/update/delete events.
We register callback functions which will be called by Informers to deliver events to controller. 

podInformer.informer().AddEventHandler(
  cache.ResourceHandler.Funcs{
    AddFunc: func(obj Interface{}) {}
    UpdateFunc: func(old, cur Interface{}) {}
    DeleteFunc: func(obj Interface{}) {}
  }
)

### Work queue
Resource Event Handler callback functions obtain an object key from the events, enqueue that key to a worker queue for further processing. Object key is combination of namespace and name of resource.

### Workers
Workers are go routines where we processes items in work queue, we can use index reference or listing wrapper to retrieve object from object key. It contains the sync handler where we run business logic of the controller.

```
func(c *ctrl) worker() {
  for c.ProcessNextItem() {}
}

func (c *ctrl) processNextItem() {
  item := c.queue.Get()
  err := c.syncHandler()
  c.handlerErr(err, item)
}
```

## Setup

### Install k8s.io/client-go in your $GOPATH. 
```
go get k8s.io/client-go/v10.0.0
```

### Install k8s.io/api-mechinery
```
go get -u k8s.io/apimachinery/v10.0.0
```

### Dependency management
If you want to use dependencies from $GOPATH, then run 'godep restore' and remove vendor folder 
```
godep restore ./v10.0.0
rm -rf ./vendor
```


### Coding the package

### imports
```
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)
```

## Connecting to API Server

### Build config
client-go provides utility function to build kubeconfig
```
config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
```
BuildConfigFromFlags builds kubeconfig from kubeconfig file path.If kubeconfig file path is empty, it falls back to inClusterConfig. If inClusterconfig fails, it fallsback to default config ~/.kube/config 

### Create clientset
clientset provides access to versioned API object.

Every resource in Kubernetes is member of API Group - core, extensions, batch, apps, etc...

Groups also contain versions, versions allow developers to introduces changes to API. Some of versions inside a group are core/v1, extensions/v1alpha1, batch/v1beta1, apps/v1beta2

NewForConfig() creates a new clientset for the given config.
```
clientset, err := kubernetes.NewForConfig(config)
```

## Listing Pods
Pods are part of core API group, with version V1. 
```
pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
```
Pod("") - Get Pods in all namespace.
List(metav1.ListOptions{}) - empty pod selector

```
fmt.Printf("No. of Pods %d\n", len(pods.Items))
```


 create a watcher for PersistentVolumeClaim resources using method Watch. 
 use the watcher to gain access to the event notifications from a Go channel via method ResultChan.


minikube addons enable heapster

kubectl get --raw "/apis/metrics.k8s.io/v1beta1/pods/memory-demo" -n mem-example

