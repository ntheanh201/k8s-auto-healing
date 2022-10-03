package handler

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	healingConfig "onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/usecase"
	"time"
)

const (
	PostgresNamespace        = "postgres"
	PostgresPoolerDeployment = "ccp-postgres-pooler"
)

func NewHandlePostgresCheckingJob(cfg *healingConfig.Config, p usecase.PostgresChecking) bool {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println(err.Error())
		return false
		//panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err.Error())
		//panic(err.Error())
		return false
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.AppsV1().Deployments(PostgresNamespace).Get(context.TODO(), PostgresPoolerDeployment, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Deployment %s not found in %s namespace\n", PostgresPoolerDeployment, PostgresNamespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting deployment %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found %s deployment in %s namespace\n", PostgresPoolerDeployment, PostgresNamespace)
		}

		_, err = p.UpsertData("test")
		if err != nil {
			return false
		}

		time.Sleep(10 * time.Minute)
	}
}
