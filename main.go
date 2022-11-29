package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	informerfactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	podinformer := informerfactory.Core().V1().Pods()
	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("Add Event Handler was called")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("Update Event handler was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete event handler was called")
		},
	})
	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podinformer.Lister().Pods("default").Get("gok8s-8fb98cf87-76wgk")
	if err != nil {
		fmt.Printf("PodInformer error: %s", err)
	}
	fmt.Println(pod)
}
