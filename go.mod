module github.com/wjam/image-registry-policy

go 1.16

require (
	github.com/containers/image v3.0.2+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/prometheus/client_golang v1.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slok/kubewebhook/v2 v2.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/tools v0.1.1
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.21.0
)
