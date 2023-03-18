package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("There are %d pods in the cluster \n", len(pods.Items))
		for i, item := range pods.Items {
			log.Printf("%d -> %s/%s", i+1, item.Namespace, item.Name)
		}
		<-time.Tick(5 * time.Second)
	}
}
