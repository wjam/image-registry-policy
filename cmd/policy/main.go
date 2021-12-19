package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	"github.com/spf13/viper"
	"github.com/wjam/image-registry-policy/internal/latest_version"
	"github.com/wjam/image-registry-policy/internal/registry"

	"net/http"
)

var configFile = flag.String("config", "./config.yaml", "config file location")

func main() {
	flag.Parse()
	err := app()
	if err != nil {
		panic(err)
	}
}

func app() error {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	registries := registry.NewValidator()

	if err := configure(*configFile, registries, log); err != nil {
		return err
	}

	handler, err := httpServer(kwhlogrus.NewLogrus(logrus.NewEntry(log)), latest_version.NewValidator(), registries)
	if err != nil {
		return err
	}

	certFile := viper.GetString("cert_file")
	keyFile := viper.GetString("key_file")
	port := fmt.Sprintf(":%s", viper.GetString("port"))

	if certFile != "" && keyFile != "" {
		log.WithField("certFile", certFile).
			WithField("keyFile", keyFile).
			WithField("port", port).
			Info("Serving HTTPS")
		if err := http.ListenAndServeTLS(port, certFile, keyFile, handler); err != nil {
			return err
		}
	} else {
		log.WithField("port", port).Info("Serving HTTP")
		if err := http.ListenAndServe(port, handler); err != nil {
			return err
		}
	}

	return nil
}
