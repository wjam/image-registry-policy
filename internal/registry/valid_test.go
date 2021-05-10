package registry

import (
	"context"
	"testing"

	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name       string
		registries []string
		images     []string
		image      string
		expected   bool
	}{
		{
			name:       "no_images_banned_registry_default_registry",
			registries: []string{"example.com"},
			image:      "busybox",
			images:     nil,
			expected:   false,
		},
		{
			name:       "no_images_banned_registry",
			registries: []string{"example.com"},
			image:      "docker.io/busybox",
			images:     nil,
			expected:   false,
		},
		{
			name:       "images_banned_registry",
			registries: []string{"example.com"},
			image:      "docker.io/busybox",
			images:     []string{"postgres:12", "example.test/dummy"},
			expected:   false,
		},
		{
			name:       "banned_registry_allowed_image",
			registries: []string{"example.com"},
			image:      "example.test/dummy:1.0",
			images:     []string{"example.test/dummy"},
			expected:   true,
		},
		{
			name:       "allowed_registry",
			registries: []string{"example.com"},
			image:      "example.com/dummy:1.0",
			images:     nil,
			expected:   true,
		},
		{
			name:       "banned_registry_allowed_image_version",
			registries: []string{"example.com"},
			image:      "example.test/dummy:1.0",
			images:     []string{"example.test/dummy:1.0"},
			expected:   true,
		},
		{
			name:       "banned_registry_different_image_version",
			registries: []string{"example.com"},
			image:      "example.test/dummy:1.1",
			images:     []string{"example.test/dummy:1.0"},
			expected:   false,
		},
		{
			name:       "banned_registry_no_image_version",
			registries: []string{"example.com"},
			image:      "example.test/dummy",
			images:     []string{"example.test/dummy:1.0"},
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			subject := NewValidator()
			subject.SetRegistries(test.registries)
			err := subject.SetImages(test.images)
			require.NoError(t, err)

			actual, err := subject.Validate(context.Background(), &model.AdmissionReview{}, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      test.name,
					Namespace: "default",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Image: test.image,
						},
					},
				},
			})
			require.NoError(t, err)

			assert.Equal(t, test.expected, actual.Valid)
		})
	}
}
