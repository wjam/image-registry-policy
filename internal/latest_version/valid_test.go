package latest_version

import (
	"context"
	"testing"

	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func TestTags(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected bool
	}{
		{
			name:     "no_tag",
			image:    "busybox",
			expected: false,
		},
		{
			name:     "latest_tag",
			image:    "busybox:latest",
			expected: false,
		},
		{
			name:     "tag",
			image:    "busybox:v1",
			expected: true,
		},
	}

	subject := NewValidator()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := subject.Validate(context.Background(), &model.AdmissionReview{}, &v1.Pod{
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
