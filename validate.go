package main

import (
	"encoding/json"
	"fmt"
	"strings"

	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	batchv1 "github.com/kubewarden/k8s-objects/api/batch/v1"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubernetes "github.com/kubewarden/policy-sdk-go/pkg/capabilities/kubernetes"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const DEPLOYMENT_KIND = "deployment"
const REPLICASET_KIND = "replicaset"
const STATEFULSET_KIND = "statefulset"
const DAEMONSET_KIND = "daemonset"
const REPLICATIONCONTROLLER_KIND = "replicationcontroller"
const CRONJOB_KIND = "cronjob"
const JOB_KIND = "job"
const POD_KIND = "pod"

func getNamespace(validationRequest kubewarden_protocol.ValidationRequest) (*corev1.Namespace, error) {

	if len(validationRequest.Request.Namespace) == 0 {
		return nil, fmt.Errorf("Admission request is missing namespace")
	}

	host := getWapcHost()

	resourceRequest := kubernetes.GetResourceRequest{
		APIVersion: "v1",
		Kind:       "Namespace",
		Name:       validationRequest.Request.Namespace,
	}

	responseBytes, err := kubernetes.GetResource(&host, resourceRequest)
	if err != nil {
		return nil, fmt.Errorf("Cannot get namespace data: %s", err)
	}
	namespace := &corev1.Namespace{}
	if err := json.Unmarshal(responseBytes, namespace); err != nil {
		return nil, fmt.Errorf("Cannot parse namespace data: %s", err)
	}
	return namespace, nil
}

func validateResourceLabels(namespaceLabels map[string]string, request kubewarden_protocol.ValidationRequest, settings Settings) ([]byte, error) {
	labelsToPropagate := make(map[string]string)
	for _, label := range settings.PropagatedLabels {
		if value, namespace_has_label := namespaceLabels[label]; namespace_has_label {
			labelsToPropagate[label] = value
		}
	}
	return updateResourceLabels(request, labelsToPropagate)

}

// propagateLabels ensures the labels defined in the meta object contains the
// same labels defined in the `labelsToPropagate` map. Returns `true` when
// the meta object has been changed
func propagateLabels(meta *metav1.ObjectMeta, labelsToPropagate map[string]string) bool {
	if meta.Labels == nil {
		meta.Labels = make(map[string]string)
	}

	hasMutation := false
	for label, newValue := range labelsToPropagate {
		if oldValue, has_label := meta.Labels[label]; !has_label || oldValue != newValue {
			meta.Labels[label] = newValue
			hasMutation = true
		}
	}
	return hasMutation
}

func updateResourceLabels(object kubewarden_protocol.ValidationRequest, labelsToPropagate map[string]string) ([]byte, error) {
	switch strings.ToLower(object.Request.Kind.Kind) {
	case DEPLOYMENT_KIND:
		deployment := appsv1.Deployment{}
		if err := json.Unmarshal(object.Request.Object, &deployment); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(deployment.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(deployment.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(deployment)
		}
	case REPLICASET_KIND:
		replicaset := appsv1.ReplicaSet{}
		if err := json.Unmarshal(object.Request.Object, &replicaset); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(replicaset.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(replicaset.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(replicaset)
		}
	case STATEFULSET_KIND:
		statefulset := appsv1.StatefulSet{}
		if err := json.Unmarshal(object.Request.Object, &statefulset); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(statefulset.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(statefulset.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(statefulset)
		}
	case DAEMONSET_KIND:
		daemonset := appsv1.DaemonSet{}
		if err := json.Unmarshal(object.Request.Object, &daemonset); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(daemonset.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(daemonset.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(daemonset)
		}
	case REPLICATIONCONTROLLER_KIND:
		replicationController := corev1.ReplicationController{}
		if err := json.Unmarshal(object.Request.Object, &replicationController); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(replicationController.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(replicationController.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(replicationController)
		}
	case CRONJOB_KIND:
		cronjob := batchv1.CronJob{}
		if err := json.Unmarshal(object.Request.Object, &cronjob); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(cronjob.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(cronjob.Spec.JobTemplate.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(cronjob)
		}
	case JOB_KIND:
		job := batchv1.Job{}
		if err := json.Unmarshal(object.Request.Object, &job); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(job.Metadata, labelsToPropagate)
		podSpecChanged := propagateLabels(job.Spec.Template.Metadata, labelsToPropagate)
		if objChanged || podSpecChanged {
			return kubewarden.MutateRequest(job)
		}
	case POD_KIND:
		pod := corev1.Pod{}
		if err := json.Unmarshal(object.Request.Object, &pod); err != nil {
			return nil, err
		}
		objChanged := propagateLabels(pod.Metadata, labelsToPropagate)
		if objChanged {
			return kubewarden.MutateRequest(pod)
		}
	default:
		return nil, fmt.Errorf("object should be one of these kinds: Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, Job, CronJob, Pod. Found %s", object.Request.Kind.Kind)
	}
	return kubewarden.AcceptRequest()
}

func validate(payload []byte) ([]byte, error) {
	// Create a ValidationRequest instance from the incoming payload
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	namespace, err := getNamespace(validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(kubewarden.Message(err.Error()), kubewarden.Code(400))
	}

	return validateResourceLabels(namespace.Metadata.Labels, validationRequest, settings)
}
