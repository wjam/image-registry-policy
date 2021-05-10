package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	"github.com/slok/kubewebhook/v2/pkg/log"
	"github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	v1 "k8s.io/api/core/v1"
)

func httpServer(log log.Logger, validators ...validating.Validator) (http.Handler, error) {
	chain := validating.NewChain(log, validators...)

	hook, err := validating.NewWebhook(validating.WebhookConfig{
		ID:        "image_check",
		Obj:       &v1.Pod{},
		Validator: chain,
		Logger:    log,
	})
	if err != nil {
		return nil, err
	}

	handler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: hook,
		Logger:  log,
	})
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/validate", handler)

	return mux, nil
}
