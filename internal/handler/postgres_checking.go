package handler

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/usecase"
	"time"
)

const (
	PostgresNamespace        = "postgres"
	PostgresPoolerDeployment = "ccp-postgres-pooler"
)

func buildConfigFromFlags(context, kubeConfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func restartPostgresDeployment(clientSet *kubernetes.Clientset) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"%s","maxSurge": "%s"}}}`, time.Now().String(), "25%", "25%")
	newDeployment, err := clientSet.AppsV1().Deployments(PostgresNamespace).Patch(context.Background(), PostgresPoolerDeployment, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	fmt.Println("new deployment: ", newDeployment)
	if err != nil {
		fmt.Printf("Error getting deployment %v\n", err)
	}
}

func showPods(clientSet *kubernetes.Clientset) {
	pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func testRestartDeployment(clientSet *kubernetes.Clientset) {
	//data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"%s","maxSurge": "%s"}}}`, time.Now().String(), "25%", "25%")
	//newDeployment, err := clientSet.AppsV1().Deployments(PostgresNamespace).Patch(context.Background(), PostgresPoolerDeployment, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	//
	//fmt.Println("new deployment: ", newDeployment)
	//if err != nil {
	//	fmt.Printf("Error getting new pooler deployment %v\n", err)
	//}
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
			//restartPostgresDeployment(clientSet)
		}

	}
	fmt.Println("DB Response: ", data)
}

func runCronJob(clientSet *kubernetes.Clientset, p usecase.PostgresChecking) {
	location, err := time.LoadLocation(config.AppConfig.App.TZ)

	if err != nil {
		fmt.Printf("Cannot get TZ from ENV")
	}

	s := gocron.NewScheduler(location)

	s.Every(2).Minute().Do(func() {
		handleUpsertData(clientSet, p)
	})

	s.StartAsync()
}

func NewHandlePostgresCheckingJob(p usecase.PostgresChecking) {
	kubeConfig := flag.String("dev-super-vcar-developer", "./config", "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	clusterContext := config.AppConfig.ClusterContext

	k8sClusterConfig, err := buildConfigFromFlags(clusterContext, *kubeConfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(k8sClusterConfig)
	if err != nil {
		log.Println(err.Error())
		return
	}

	runCronJob(clientSet, p)
}
