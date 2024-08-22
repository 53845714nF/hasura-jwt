# Kubernetes Deployment

## Prerequisites

- Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- On Windows, install
  [Docker Desktop](https://www.docker.com/products/docker-desktop)
  an enable Kubernetes
- On Linux, install [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
- Install [Helm](https://helm.sh/docs/intro/install/)

## Local development

Set the following records in your `etc/hosts`:

```text
127.0.0.1   hasura.docker.internal
127.0.0.1   nginx.local
127.0.0.1   hasura-jwt.docker.internal
```

## Deploy to Kubernetes

Use `make ingress` to install nginx ingress on your Cluster.
If you don't already have it.

Use `make all` to create this app.

## Uninstall

Use `make clean` to delete the app from the Kubernetes Cluster.
