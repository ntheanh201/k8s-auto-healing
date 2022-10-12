package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/notify/webhook"
	prommodel "github.com/prometheus/common/model"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"onroad-k8s-auto-healing/internal/entity"
	"onroad-k8s-auto-healing/internal/handler"
)

const (
	SeverityCritical                 = "critical"
	StuckThingsboardDataRateLimitter = "Stuck Thingsboard data rate limitter"
)

type alertGroupV4 webhook.Message

func (a alertGroupV4) toDomain() (*entity.AlertGroup, error) {
	//if a.Version != "4" {
	//	return nil, errors.New("not supported alert group version")
	//}

	// Map alerts.
	alerts := make([]entity.Alert, 0, len(a.Alerts))
	for _, alert := range a.Alerts {
		modelAlert := entity.Alert{
			ID:           alert.Fingerprint,
			Name:         alert.Labels[prommodel.AlertNameLabel],
			StartsAt:     alert.StartsAt,
			EndsAt:       alert.EndsAt,
			Status:       alertStatusToDomain(alert.Status),
			Labels:       alert.Labels,
			Annotations:  alert.Annotations,
			GeneratorURL: alert.GeneratorURL,
		}
		alerts = append(alerts, modelAlert)
	}

	ag := &entity.AlertGroup{
		ID:     a.GroupKey,
		Labels: a.GroupLabels,
		Alerts: alerts,
	}

	return ag, nil
}

func alertStatusToDomain(st string) entity.AlertStatus {
	switch prommodel.AlertStatus(st) {
	case prommodel.AlertFiring:
		return entity.AlertStatusFiring
	case prommodel.AlertResolved:
		return entity.AlertStatusResolved
	default:
		return entity.AlertStatusUnknown
	}
}

type webhookRoutes struct {
	clientSet *kubernetes.Clientset
}

func NewWebhookRoutes(handler *gin.RouterGroup, clientSet *kubernetes.Clientset) {
	w := &webhookRoutes{
		clientSet: clientSet,
	}
	h := handler.Group("/webhooks")
	{
		h.POST("/prometheus", w.handlePrometheus)
	}
}

func (w *webhookRoutes) handlePrometheus(ctx *gin.Context) {
	reqAlerts := alertGroupV4{}
	err := ctx.BindJSON(&reqAlerts)
	if err != nil {
		log.Printf("error unmarshalling JSON: %s\n", err)
		_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	model, err := reqAlerts.toDomain()
	if err != nil {
		log.Printf("error mapping to domain models: %s\n", err)
		_ = ctx.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	for _, alert := range model.Alerts {
		severity := alert.Labels["severity"]
		namespace := alert.Labels["namespace"]
		pod := alert.Labels["pod"]
		alertName := alert.Name

		if severity == SeverityCritical {
			switch alertName {
			case StuckThingsboardDataRateLimitter:
				handler.NewTelemetryHandler(w.clientSet, pod, namespace)
			}
		}

	}

	ctx.JSON(http.StatusOK, "OK")
}
