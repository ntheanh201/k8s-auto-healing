package handler

import (
	"context"
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"onroad-k8s-auto-healing/config"
)

type ClusterClient struct {
	ClientSet *kubernetes.Clientset
}

func buildConfigFromFlags(context, kubeConfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
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

func NewClientSetCluster() *ClusterClient {
	kubeConfig := flag.String("dev-super-vcar-developer", "./kube-config", "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	clusterContext := config.AppConfig.ClusterContext

	k8sClusterConfig, err := buildConfigFromFlags(clusterContext, *kubeConfig)
	if err != nil {
		//panic(err)
		log.Println(err.Error())
		log.Println("Cannot get k8s cluster config")
		return nil
	}

	clientSet, err := kubernetes.NewForConfig(k8sClusterConfig)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &ClusterClient{ClientSet: clientSet}
}
