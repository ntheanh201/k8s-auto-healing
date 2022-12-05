package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	MonitoringNamespace = "monitoring"
	FluentBitDaemonSet  = "fluent-bit"
)

type FluentBitHandler struct {
	clusterClient *ClusterClient
}

func restartDaemonSetFluentBit(clientSet *kubernetes.Clientset) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"%s","maxSurge": "%s"}}}`, time.Now().String(), "25%", "25%")
	newDaemonSet, err := clientSet.AppsV1().DaemonSets(MonitoringNamespace).Patch(context.Background(), FluentBitDaemonSet, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	fmt.Println("new daemonSet: ", newDaemonSet)
	if err != nil {
		fmt.Printf("Error getting daemonSet %v\n", err)
	}
}

func handleRestartFluentBit(clientSet *kubernetes.Clientset) {
	_, err := clientSet.AppsV1().DaemonSets(MonitoringNamespace).Get(context.TODO(), FluentBitDaemonSet, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("DaemonSet %s not found in %s namespace\n", FluentBitDaemonSet, MonitoringNamespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting DaemonSet %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found %s DaemonSet in %s namespace\n", FluentBitDaemonSet, MonitoringNamespace)
		restartDaemonSetFluentBit(clientSet)
	}
}

func NewFluentBitHandler(client *ClusterClient) *FluentBitHandler {
	return &FluentBitHandler{clusterClient: client}
}

func (f *FluentBitHandler) StartNewJob() {
	s := gocron.NewScheduler(time.UTC)

	// run every 2 hours
	s.Every(2).Hour().Do(func() {
		handleRestartFluentBit(f.clusterClient.ClientSet)
	})

	s.StartAsync()
}

//func (c *ClusterClient) NewFluentBitHandler() {
//	s := gocron.NewScheduler(time.UTC)
//
//	// run every 2 hours
//	s.Every(2).Hour().Do(func() {
//		handleRestartFluentBit(c.ClientSet)
//	})
//
//	s.StartAsync()
//}
