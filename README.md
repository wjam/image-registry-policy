# image-registry-policy

This is a [ValidatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
that helps to enforce restrictions around where images can be pulled from.

It allows allowlisting whole registries as well as individual images and specific tags of images. The difference 
between `docker.io/library/postgres:latest`, `docker.io/library/postgres`, `library/postgres` & `postgres` is handled.
It also blocks the use of the `latest` tag - either explicitly set or through no tag set.

The configuration file looks like this and will be automatically reloaded when it changes:
```yaml
log_level: INFO
allowed_registries:
  - 602401143452.dkr.ecr.eu-west-1.amazonaws.com
  - quay.io
allowed_images:
  - docker.io/library/postgres:12
  - library/golang:1.16
  - vault:1.7.1
```
