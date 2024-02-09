# Kubernetes Deployment 

## Prerequisites
  - Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
  - On Windows, install [Docker Desktop](https://www.docker.com/products/docker-desktop) an enable Kubernetes
  - On Linux, install [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
  - Install [Helm](https://helm.sh/docs/intro/install/)

## Deploy to Kubernetes

```bash
helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --create-namespace
helm repo add hasura https://hasura.github.io/helm-charts
helm repo update
helm install hasura -f ./values.yaml hasura/graphql-engine
```

## Uninstall

```bash
helm uninstall hasura
```