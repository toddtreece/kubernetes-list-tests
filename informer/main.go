package main

import (
	"flag"
	"fmt"
	"time"

	clientset "github.com/toddtreece/kubernetes-list-tests/informer/pkg/client/clientset/versioned"
	informers "github.com/toddtreece/kubernetes-list-tests/informer/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig.")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(client, 10*time.Minute)

	dashboardInformer := informerFactory.Dashboard().V1alpha1().Dashboards()

	dashboardInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Added: ", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Updated: ", newObj)
		},
	})

	stop := make(chan struct{})
	defer close(stop)

	informerFactory.Start(stop)

	select {}
}
