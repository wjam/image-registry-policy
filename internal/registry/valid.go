package registry

import (
	"context"
	"fmt"

	"github.com/containers/image/docker/reference"
	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ validating.Validator = &Valid{}

func NewValidator() *Valid {
	return &Valid{}
}

type Valid struct {
	registries []string
	images     []reference.Named
}

func (v *Valid) SetRegistries(registries []string) {
	v.registries = registries
}

func (v *Valid) SetImages(images []string) error {
	var parsed []reference.Named

	for _, image := range images {
		named, err := reference.ParseNormalizedNamed(image)
		if err != nil {
			return err
		}
		parsed = append(parsed, named)
	}

	v.images = parsed

	return nil
}

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

		registry := reference.Domain(image)

		if !contains(registry, v.registries) && !v.allowedImage(image) {
			return &validating.ValidatorResult{
				Valid:   false,
				Message: fmt.Sprintf("Registry %s is not allowed", registry),
			}, nil
		}
	}

	return &validating.ValidatorResult{
		Valid:   true,
		Message: "all registries valid",
	}, nil
}

func (v *Valid) allowedImage(image reference.Named) bool {
	for _, allowed := range v.images {
		if subset(allowed, image) {
			return true
		}
	}

	return false
}

func contains(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func subset(super reference.Named, sub reference.Named) bool {
	if reference.Domain(super) != reference.Domain(sub) {
		return false
	}

	if super.Name() != sub.Name() {
		return false
	}

	if superTag, ok := super.(reference.Tagged); ok {
		if subTag, ok := sub.(reference.Tagged); ok {
			if superTag != subTag {
				return false
			}
		} else {
			return false
		}
	}

	return true
}
