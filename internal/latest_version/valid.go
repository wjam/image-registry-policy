package latest_version

import (
	"context"
	"fmt"

	"github.com/containers/image/docker/reference"
	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewValidator() validating.Validator {
	return &Valid{}
}

type Valid struct{}

func (v *Valid) Validate(_ context.Context, _ *model.AdmissionReview, obj metav1.Object) (*validating.ValidatorResult, error) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return nil, fmt.Errorf("not a pod")
	}

	for _, container := range pod.Spec.Containers {

		image, err := reference.ParseNormalizedNamed(container.Image)
		if err != nil {
			return nil, err
		}

		image = reference.TagNameOnly(image)

		if v, ok := image.(reference.Tagged); ok && v.Tag() == "latest" {
			return &validating.ValidatorResult{
				Valid:   false,
				Message: "Latest version is not allowed",
			}, nil
		}
	}

	return &validating.ValidatorResult{
		Valid:   true,
		Message: "all images use tags",
	}, nil
}
