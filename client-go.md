
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
clientset holds the Kubernetes client interface to interact with API server.

Every resource in Kubernetes is a member of API Group - core, extensions, batch, apps, etc...

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


