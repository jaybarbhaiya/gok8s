package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
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
	// config.Timeout = 120 * time.Second

	// clientset can be used to interactive with the k8s resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error while creating the client set: %s", err.Error())
	}

	// get the pods in the default name space
	ctx := context.Background()
	namespace := "default"

	podsClient := clientset.CoreV1().Pods(namespace)
	// pod, err := podsClient.Get(ctx, "gok8s-8fb98cf87-swctq", metav1.GetOptions{})
	// if err != nil {
	// 	fmt.Printf("Failed to get pod: %s", err.Error())
	// }
	// fmt.Printf("pod type: %T\n", pod)

	podList, podListErr := podsClient.List(ctx, metav1.ListOptions{})
	if podListErr != nil {
		fmt.Printf("Failed to get PodList: %s", podListErr.Error())
	}

	var pod v1.Pod

	for _, podItem := range podList.Items {
		if strings.Contains(podItem.Name, "gok8s") {
			pod = podItem
		}
	}

	if pod.Name == "" {
		panic("Relevant Pod not found in the podList")
	}
	fmt.Printf("Relevant Pod found: %s", pod.Name)

	if pod.ObjectMeta.Labels["jblabel"] != "jblabelvalue" {
		pod.ObjectMeta.Labels["jblabel"] = "jblabelvalue"
	}

	_, updateErr := podsClient.Update(ctx, &pod, metav1.UpdateOptions{})
	if updateErr != nil {
		fmt.Printf("Failed to update pod: %s", updateErr.Error())
	}

	for labelKey, labelValue := range pod.ObjectMeta.Labels {
		fmt.Printf("%v: %v\n", labelKey, labelValue)
	}

}
