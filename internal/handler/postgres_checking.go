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

	"onroad-k8s-auto-healing/internal/usecase"
)

const (
	PostgresNamespace        = "postgres"
	PostgresPoolerDeployment = "ccp-postgres-pooler"
)

func restartPostgresDeployment(clientSet *kubernetes.Clientset) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"%s","maxSurge": "%s"}}}`, time.Now().String(), "25%", "25%")
	newDeployment, err := clientSet.AppsV1().Deployments(PostgresNamespace).Patch(context.Background(), PostgresPoolerDeployment, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	fmt.Println("new deployment: ", newDeployment)
	if err != nil {
		fmt.Printf("Error getting deployment %v\n", err)
	}
}

func handleUpsertData(clientSet *kubernetes.Clientset, p usecase.PostgresChecking) {
	t := time.Now()
	fmt.Println("Upserting new checking data")

	data, err := p.UpsertCheckingData(t.Format("2006-01-02 15:04:05"))

	if err != nil {
		_, err = clientSet.AppsV1().Deployments(PostgresNamespace).Get(context.TODO(), PostgresPoolerDeployment, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Deployment %s not found in %s namespace\n", PostgresPoolerDeployment, PostgresNamespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting deployment %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found %s deployment in %s namespace\n", PostgresPoolerDeployment, PostgresNamespace)
			restartPostgresDeployment(clientSet)
		}

	}
	fmt.Println("DB Response: ", data)
}

func (c *ClusterClient) NewHandlePostgresCheckingJob(p usecase.PostgresChecking) {
	s := gocron.NewScheduler(time.UTC)

	// run every 5 minutes
	s.Every(5).Minute().Do(func() {
		handleUpsertData(c.ClientSet, p)
	})

	s.StartAsync()
}
