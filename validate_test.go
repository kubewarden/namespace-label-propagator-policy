package main

import (
	"encoding/json"
	"fmt"

	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	batchv1 "github.com/kubewarden/k8s-objects/api/batch/v1"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	capabilities "github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/mailru/easyjson"
	"testing"
)

const SHOULD_ACCEPT = true
const SHOULD_REJECT = false
const SHOULD_MUTATE = true
const NO_MUTATION = false
const TEST_NAMESPACE = "default"

func buildValidationRequest(propagatedLabels []string, resource easyjson.Marshaler, kind string) ([]byte, error) {
	settings := Settings{PropagatedLabels: propagatedLabels}
	payload, err := kubewarden_testing.BuildValidationRequest(resource, &settings)

	if err != nil {
		return nil, err
	}
	payload, err = updateValidationRequestKindAndNamespace(payload, kind)
	if err != nil {
		return nil, err
	}
	return payload, nil

}
func buildWapcClient(namespaceLabels map[string]string) error {
	clientResponse := corev1.Namespace{
		Metadata: &metav1.ObjectMeta{
			Labels: namespaceLabels,
		},
	}
	var err error
	wapcClient, err = capabilities.NewSuccessfulMockWapcClient(clientResponse)
	if err != nil {
		return err
	}
	return nil

}

func basicResposeValidation(responsePayload []byte, accepted, should_mutate bool) (*kubewarden_protocol.ValidationResponse, error) {
	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		return nil, fmt.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != accepted {
		return nil, fmt.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, *response.Code)
	}

	if response.MutatedObject == nil && should_mutate {
		return nil, fmt.Errorf("Missing mutated resource")
	}
	return &response, nil

}

func validateLabels(resourceLabels, expectedLabels map[string]string) error {
	for expectedLabel, expectedValue := range expectedLabels {
		if resourceValue, found := resourceLabels[expectedLabel]; found {
			if resourceValue != expectedValue {
				return fmt.Errorf("Resource label \"%s\" expected value:  \"%s\". Found \"%s\"", expectedLabel, expectedValue, resourceValue)
			}
		} else {
			return fmt.Errorf("Mutated resource missing label \"%s\"", expectedLabel)
		}

	}

	if len(resourceLabels) != len(expectedLabels) {
		return fmt.Errorf("Mutated resource contains %d labels. But the expected is %d", len(resourceLabels), len(expectedLabels))
	}
	return nil
}

func updateValidationRequestKindAndNamespace(payload []byte, kind string) ([]byte, error) {
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := easyjson.Unmarshal(payload, &validationRequest)
	if err != nil {
		return nil, err
	}
	validationRequest.Request.Kind.Kind = kind
	validationRequest.Request.Namespace = TEST_NAMESPACE
	return easyjson.Marshal(validationRequest)
}

func TestPodWithNoLabels(t *testing.T) {
	propagatedLabels := []string{"testing"}
	namespaceLabels := map[string]string{
		"testing":  "foo",
		"testing2": "zpto",
	}
	expectedLabels := map[string]string{
		"testing": "foo",
	}

	resource := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}

	payload, err := buildValidationRequest(propagatedLabels, resource, POD_KIND)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	err = buildWapcClient(namespaceLabels)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	response, err := basicResposeValidation(responsePayload, SHOULD_ACCEPT, SHOULD_MUTATE)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	mutatedResourceJSON, err := json.Marshal(response.MutatedObject.(map[string]interface{}))
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if err := validateLabels(resource.Metadata.Labels, expectedLabels); err != nil {
		t.Error(err.Error())
	}
}

func TestPodLabelsShouldNotMutateWithItHasTheExpectedValue(t *testing.T) {
	propagatedLabels := []string{"testing"}
	namespaceLabels := map[string]string{
		"testing":  "foo",
		"testing2": "zpto",
	}

	resource := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"testing":  "foo",
				"testing2": "zzz",
			},
		},
	}

	payload, err := buildValidationRequest(propagatedLabels, resource, POD_KIND)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	err = buildWapcClient(namespaceLabels)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	_, err = basicResposeValidation(responsePayload, SHOULD_ACCEPT, NO_MUTATION)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestLabelsShouldOverwrittenLabelsOnlyDefinedInSettings(t *testing.T) {
	cases := []struct {
		propagatedLabels []string
		namespaceLabels  map[string]string
		expectedLabels   map[string]string
		resource         easyjson.Marshaler
		kind             string
		accept           bool
		mutate           bool
	}{
		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo"},
			corev1.Pod{Metadata: &metav1.ObjectMeta{Name: "test", Namespace: "default"}},
			POD_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},

		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			batchv1.Job{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &batchv1.JobSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			JOB_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			batchv1.CronJob{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &batchv1.CronJobSpec{
					JobTemplate: &batchv1.JobTemplateSpec{
						Spec: &batchv1.JobSpec{
							Template: &corev1.PodTemplateSpec{
								Metadata: &metav1.ObjectMeta{
									Name:   "podtest",
									Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
								},
							},
						},
					},
				},
			},
			CRONJOB_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},

		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			corev1.ReplicationController{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &corev1.ReplicationControllerSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			REPLICATIONCONTROLLER_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			appsv1.DaemonSet{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &appsv1.DaemonSetSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			DAEMONSET_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			appsv1.StatefulSet{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &appsv1.StatefulSetSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			STATEFULSET_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
		{

			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			appsv1.ReplicaSet{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &appsv1.ReplicaSetSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			REPLICASET_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
		{
			[]string{"testing"},
			map[string]string{"testing": "foo", "testing2": "zpto"},
			map[string]string{"testing": "foo", "testing2": "zzz"},
			appsv1.Deployment{
				Metadata: &metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Labels:    map[string]string{"testing": "bar", "testing2": "zzz"},
				},
				Spec: &appsv1.DeploymentSpec{
					Template: &corev1.PodTemplateSpec{
						Metadata: &metav1.ObjectMeta{
							Name:   "podtest",
							Labels: map[string]string{"testing": "pod-bar", "testing2": "zzz"},
						},
					},
				},
			},
			DEPLOYMENT_KIND,
			SHOULD_ACCEPT, SHOULD_MUTATE,
		},
	}

	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			payload, err := buildValidationRequest(tc.propagatedLabels, tc.resource, tc.kind)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			err = buildWapcClient(tc.namespaceLabels)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			responsePayload, err := validate(payload)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			response, err := basicResposeValidation(responsePayload, tc.accept, tc.mutate)
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			mutatedResourceJSON, err := json.Marshal(response.MutatedObject.(map[string]interface{}))
			if err != nil {
				t.Errorf("Unexpected error: %+v", err)
			}

			switch tc.kind {
			case POD_KIND:
				resource := corev1.Pod{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}

				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case DEPLOYMENT_KIND:
				resource := appsv1.Deployment{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case REPLICASET_KIND:
				resource := appsv1.ReplicaSet{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case DAEMONSET_KIND:
				resource := appsv1.DaemonSet{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case STATEFULSET_KIND:
				resource := appsv1.StatefulSet{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case REPLICATIONCONTROLLER_KIND:
				resource := corev1.ReplicationController{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case JOB_KIND:
				resource := batchv1.Job{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			case CRONJOB_KIND:
				resource := batchv1.CronJob{}
				if err := easyjson.Unmarshal(mutatedResourceJSON, &resource); err != nil {
					t.Errorf("Unexpected error: %+v", err)
				}
				if err := validateLabels(resource.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
				if err := validateLabels(resource.Spec.JobTemplate.Spec.Template.Metadata.Labels, tc.expectedLabels); err != nil {
					t.Error(err.Error())
				}
			default:
				t.Errorf("Unexpected kind: %s", tc.kind)
			}
		})
	}
}
