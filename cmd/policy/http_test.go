package main

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:embed testdata/admin_review.json
var testAdminReview string

func TestHTTPServer_healthz(t *testing.T) {
	subject := newHttpServer(t)

	assert.HTTPSuccess(t, subject.ServeHTTP, "GET", "/healthz", nil)
}

func TestHTTPServer_metrics(t *testing.T) {
	subject := newHttpServer(t)

	assert.HTTPSuccess(t, subject.ServeHTTP, "GET", "/metrics", nil)
}

func TestHTTPServer_validate(t *testing.T) {
	subject := newHttpServer(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/validate", strings.NewReader(testAdminReview))
	require.NoError(t, err)
	subject.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func newHttpServer(t *testing.T) http.Handler {
	subject, err := httpServer(kwhlogrus.NewLogrus(logrus.NewEntry(logrus.New())), &dummyValidator{})
	require.NoError(t, err)
	return subject
}

var _ validating.Validator = &dummyValidator{}

type dummyValidator struct {
}

func (d *dummyValidator) Validate(context.Context, *model.AdmissionReview, metav1.Object) (*validating.ValidatorResult, error) {
	return &validating.ValidatorResult{Valid: true}, nil
}
