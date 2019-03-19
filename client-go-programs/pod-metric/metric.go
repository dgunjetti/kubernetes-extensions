package main

import (
	"fmt"

	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	var kubeconfig *string

	// flag for kubeconfig
	kubeconfig = flag.String("kubeconfig",
		"",
		"absolute path to the kubeconfig file")

	flag.Parse()

	fmt.Println(*kubeconfig)

	// Build configuration kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	listOptions := metav1.ListOptions{
		LabelSelector: "run=resource-consumer",
	}

	podMetrics, err := mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(listOptions)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, podMetric := range podMetrics.Items {
		podContainers := podMetric.Containers
		for _, container := range podContainers {
			cpuQuantity, ok := container.Usage.Cpu().AsInt64()
			memQuantity, ok := container.Usage.Memory().AsInt64()
			if !ok {
				return
			}
			msg := fmt.Sprintf("Container Name: %s \n CPU usage: %d \n Memory usage: %d", container.Name, cpuQuantity, memQuantity)
			fmt.Println(msg)
		}
	}
}
