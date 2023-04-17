package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var hazelcastlog = logf.Log.WithName("hazelcast-resource")

func (r *Hazelcast) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-hazelcast-com-v1beta1-hazelcast,mutating=false,failurePolicy=ignore,sideEffects=None,groups=hazelcast.com,resources=hazelcasts,verbs=create;update,versions=v1beta1,name=vvhazelcast.kb.io,admissionReviewVersions=v1
// Role related to webhooks

var _ webhook.Validator = &Hazelcast{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Hazelcast) ValidateCreate() error {
	hazelcastlog.Info("validate create", "name", r.Name)
	return ValidateHazelcastSpec(r)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Hazelcast) ValidateUpdate(old runtime.Object) error {
	hazelcastlog.Info("validate update", "name", r.Name)
	return ValidateHazelcastSpec(r)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Hazelcast) ValidateDelete() error {
	hazelcastlog.Info("validate delete", "name", r.Name)
	return nil
}