# Client Libraries
Kubernetes provides client libraries to build application over Kubernetes REST API's.
Here we will be looking into client-go and api-machinery.

## api-machinery
Provides scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes like API resources.

## client-go
Client-go provides packages and utilities to access Kubernetes API resources. It handles the common tasks such as authentication. Discover and use service account to authenticate with API server if the client is running inside the cluster. For the clients running outside the cluster it can read kubeconfig file to get credentials and API server address.

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

kubernetes package contains the clientset to access Kubernetes API.
tools/clientcmd package contains utilities to build client from kubeconfig file.

## Connecting to API Server

### Build config
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

List(metav1.ListOptions{}) - Get all Pods, we are specifying empty Pod selector

```
fmt.Printf("No. of Pods %d\n", len(pods.Items))
```

```
fmt.Println("Listing Pods..")
for _, p := range pods.Items {
	fmt.Printf("%s\n", p.ObjectMeta.Name)
}
```


