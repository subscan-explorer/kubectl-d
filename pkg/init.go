package pkg

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func InitK8sClient(kubeconfigs ...string) *kubernetes.Clientset {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if len(kubeconfigs) != 0 && kubeconfigs[0] != "" {
		kubeconfig = kubeconfigs[0]
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
