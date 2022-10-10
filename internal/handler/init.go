package handler

import (
	"context"
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"onroad-k8s-auto-healing/config"
	"onroad-k8s-auto-healing/internal/entity"
)

type Module struct {
	Db *dbEntity
}

type dbEntity struct {
	conn                *gorm.DB
	PostgresCheckingOrm entity.CheckEntityOrm
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

func NewClientSetCluster() *kubernetes.Clientset {
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
	return clientSet
}

func NewHandler() (module *Module, err error) {
	// Initialize DB
	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s statement_cache_mode=describe",
			config.AppConfig.Db.Host, config.AppConfig.Db.Port, config.AppConfig.Db.Database,
			config.AppConfig.Db.Username, config.AppConfig.Db.Password),
	), &gorm.Config{})

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(0)

	if err != nil {
		log.Println("[INIT] failed connecting to PostgreSQL")
		return
	}
	log.Println("[INIT] connected to PostgreSQL")

	// Compose handler modules
	return &Module{
		Db: &dbEntity{
			conn:                db,
			PostgresCheckingOrm: entity.NewPostgresCheckingOrm(db),
		},
	}, nil

}
