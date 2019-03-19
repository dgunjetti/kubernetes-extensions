package main

import (
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
)

type PodStatusCheck struct {
	PodName       string
	ContainerName string
	PodCPU        string
	PodMemory     string
}

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

	// Create a client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Get all pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("No. of Pods in cluster %d\n", len(pods.Items))

	fmt.Println("Listing Pods..")
	for _, p := range pods.Items {
		fmt.Printf("%s\n", p.ObjectMeta.Name)
	}

	// Create deployment
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	size := int32(2)
	dep := makeNginxDeployment(size)
	result, err := deploymentsClient.Create(dep)
	if err != nil {
		fmt.Println("failed to create nginx deployment")
		os.Exit(1)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	listOptions := metav1.ListOptions{
		LabelSelector: "app=nginx",
	}
	watcher, err := deploymentsClient.Watch(metav1.ListOptions{})
	if err != nil {
		fmt.Println("failed to set watch on deployment ", err)
		os.Exit(1)
	}

	ch := watcher.ResultChan()
	for event := range ch {
		fmt.Println("got event", event.Type)

		if event.Type == watch.Modified {
			fmt.Println("got modify event")
			depOut, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				fmt.Println("unexpected type")
			}
			fmt.Printf("Ready replicas %d\n", depOut.Status.ReadyReplicas)
			if depOut.Status.ReadyReplicas == size {
				fmt.Println("replicas are ready")
				autoScale, err := deploymentsClient.GetScale(result.GetObjectMeta().GetName(),
					metav1.GetOptions{})
				if err != nil {
					fmt.Println("failed to get scale of nginx deployment")
					os.Exit(1)
				}
				autoScale.Spec.Replicas = 3
				_, err = deploymentsClient.UpdateScale(result.GetObjectMeta().GetName(),
					autoScale)
				if err != nil {
					fmt.Println("failed to update nginx deployment")
					os.Exit(1)
				}
				//	podMetricsList, _ := mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(metav1.ListOptions{})
				//podMetricsList, _ := mc.MetricsV1beta1().PodMetricses(apiv1.NamespaceDefault).List(listOptions)
				podMetricsList, err := mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(listOptions)

				for _, pm := range podMetricsList.Items {
					for _, c := range pm.Containers {
						sc := PodStatusCheck{
							ContainerName: c.Name,
							PodName:       pm.Name,
						}
						sc.PodCPU = c.Usage.Cpu().String()
						sc.PodMemory = c.Usage.Memory().String()
						fmt.Println("metrics")
						fmt.Println(sc)
					}
				}
				break
			}
		}
	}
}

// makeNginxDeployment creates nginx deployment with 2 replicas of latest nginx image.
func makeNginxDeployment(size int32) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "nginx",
							Image: "dgunjetti/nginx",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							LivenessProbe: &apiv1.Probe{
								InitialDelaySeconds: 3,
								PeriodSeconds:       15,
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{IntVal: 80},
										HTTPHeaders: []apiv1.HTTPHeader{{
											Name:  "Custom-Header",
											Value: "Awesome",
										}},
									},
								},
							},
							ReadinessProbe: &apiv1.Probe{
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{IntVal: 80},
										HTTPHeaders: []apiv1.HTTPHeader{{
											Name:  "Custom-Header",
											Value: "Awesome",
										}},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func int32Ptr(i int32) *int32 { return &i }
