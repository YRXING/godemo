package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var (
	label  map[string]string
	PodIPs []string
)

func main() {
	//label := make(map[string]string, 1)
	//label["app"]="calico-godemo"
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig",
			filepath.Join(home, ".kube", "config"), "optional absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	//uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	//creates the clientset
	dynamicClient, err := dynamic.NewForConfig(config)
	//clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unstructObj, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.TODO(), metav1.ListOptions{Limit: 500})
	if unstructObj == nil {
		fmt.Printf("no resource")
		os.Exit(1)
	}
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}
	for _, d := range podList.Items {
		fmt.Printf("namespace:%v \t name:%v \t status:%+v\n", d.Namespace, d.Name, d.Status.Phase)
	}
	fmt.Println("Finished executing")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

//kubeclient, err := kubernetes.NewForConfig(&rest.Config{
//Host: "http://localhost:10550",
//})
//if err != nil {
//fmt.Println("failed to new a kubeclient for testing")
//}
//podlist,err :=kubeclient.CoreV1().Pods("touchuyht").List(context.TODO(),metav1.ListOptions{})
////ep, err := kubeclient.CoreV1().Endpoints("default").Get(context.TODO(), "kubernetes", metav1.GetOptions{})
//if err != nil {
//fmt.Println("failed to get, %v", err)
//}
//for _,pod :=  range podlist.Items{
//fmt.Printf("namespace:%v \t name:%v \t status:%+v\n", pod.Namespace,pod.Name,pod.Status.Phase)
//}
