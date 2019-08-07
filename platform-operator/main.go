package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientset "github.com/cloud-ark/kubeplus/platform-operator/pkg/client/clientset/versioned"
	informers "github.com/cloud-ark/kubeplus/platform-operator/pkg/client/informers/externalversions"
	"github.com/cloud-ark/kubeplus/platform-operator/pkg/signals"

	"k8s.io/client-go/rest"
)

/*
var (
	masterURL  string
	kubeconfig string
)
*/

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	// creates the in-cluster config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	//cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	//if err != nil {
	//	glog.Fatalf("Error building kubeconfig: %s", err.Error())
	//}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	platformOperatorClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	platformInformerFactory := informers.NewSharedInformerFactory(platformOperatorClient, time.Second*30)

	controller := NewController(kubeClient, platformOperatorClient, kubeInformerFactory, platformInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go platformInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

/*
func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
*/

