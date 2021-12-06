module github.com/wjam/image-registry-policy

go 1.16

require (
	github.com/containers/image v3.0.2+incompatible
	github.com/fsnotify/fsnotify v1.5.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slok/kubewebhook/v2 v2.2.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
)
