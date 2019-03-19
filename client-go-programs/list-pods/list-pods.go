// list-pods.go
// go run list-pods.go --kubeconfig=/Users/dgunjetti/.kube/kind-config-kind

package main

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("No. of Pods in cluster %d\n", len(pods.Items))

	fmt.Println("Listing Pods..")
	for _, p := range pods.Items {
		fmt.Printf("%s\n", p.ObjectMeta.Name)
	}
}
