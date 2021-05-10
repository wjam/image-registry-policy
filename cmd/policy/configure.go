package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func configure(file string, registries imagesAndRegistries, log *logrus.Logger) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(file)

	viper.SetDefault("port", "8000")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	updateConfig(registries, log)

	viper.WatchConfig()
	viper.OnConfigChange(func(_ fsnotify.Event) {
		updateConfig(registries, log)
	})

	return nil
}

func updateConfig(v imagesAndRegistries, log *logrus.Logger) {
	logLevel := viper.GetString("log_level")
	registries := viper.GetStringSlice("allowed_registries")
	images := viper.GetStringSlice("allowed_images")

	log.WithField("logLevel", logLevel).
		WithField("registries", registries).
		WithField("images", images).
		Info("Loaded config")

	if logLevel != "" {
		if level, err := logrus.ParseLevel(logLevel); err == nil {
			log.Level = level
		}
	}
	v.SetRegistries(registries)
	if err := v.SetImages(images); err != nil {
		log.Errorf("images %s is invalid: %s", images, err)
	}
}

type imagesAndRegistries interface {
	SetImages([]string) error
	SetRegistries([]string)
}
