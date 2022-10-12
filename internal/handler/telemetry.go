package handler

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func handleDeleteThingsboardRuleEnginePod(clientSet *kubernetes.Clientset, podName, namespace string) {
	_, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s not found in %s namespace\n", podName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found %s Pod in %s namespace\n", podName, namespace)
		err = clientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewTelemetryHandler(clientSet *kubernetes.Clientset, podName, namespace string) {
	handleDeleteThingsboardRuleEnginePod(clientSet, podName, namespace)
}
