package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/jay/.kube/config", "Location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error while Building config from Flags: %s", err.Error())

		// the below code will be triggered with there is an error and will look for the InClusterConfig
		// since the hardcoded path of the config file does not exist in the cluster
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error in getting inClusterConfig: %s", err.Error())
		}
	}

	// clientset can be used to interactive with the k8s resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error while creating the client set: %s", err.Error())
	}

	ctx := context.Background()
	namespace := "default"
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods from %s namespace: %s", namespace, err.Error())
	}

	fmt.Printf("Pods from %s Namespace:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing desployments from %s namespace: %s", namespace, err.Error())
	}
	fmt.Printf("Deployments from %s Namespace:\n", namespace)
	for _, deployment := range deployments.Items {
		fmt.Printf("%s\n", deployment.Name)
	}
}