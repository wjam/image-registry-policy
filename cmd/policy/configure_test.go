package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigure_reloadsConfiguration(t *testing.T) {
	validator := &dummyImageAndRegistries{}
	log := logrus.New()
	dir := t.TempDir()

	file := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(file, []byte(`
log_level: INFO
allowed_registries:
  - example.org
  - bbc.co.uk
`), 0600)
	require.NoError(t, err)

	err = configure(file, validator, log)
	require.NoError(t, err)

	assert.Equal(t, logrus.InfoLevel, log.Level)
	assert.Equal(t, []string{"example.org", "bbc.co.uk"}, validator.registries)

	err = os.WriteFile(file, []byte(`
log_level: DEBUG
allowed_registries:
  - example.test
allowed_images:
  - foo
  - bar
`), 0600)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return log.Level == logrus.DebugLevel
	}, time.Second*1, time.Millisecond*100)

	assert.Equal(t, []string{"example.test"}, validator.registries)
	assert.Equal(t, []string{"foo", "bar"}, validator.images)
}

var _ imagesAndRegistries = &dummyImageAndRegistries{}

type dummyImageAndRegistries struct {
	images     []string
	registries []string
}

func (d *dummyImageAndRegistries) SetImages(images []string) error {
	d.images = images
	return nil
}

func (d *dummyImageAndRegistries) SetRegistries(registries []string) {
	d.registries = registries
}
