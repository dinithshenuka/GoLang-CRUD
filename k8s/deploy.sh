#!/bin/bash

# Exit on error
set -e

echo "Deploying GoLang-CRUD to Minikube..."

# Start Minikube if not already running
if ! minikube status > /dev/null 2>&1; then
  echo "Starting Minikube..."
  minikube start
fi

# Point docker to minikube's daemon
echo "Configuring Docker to use Minikube's daemon..."
eval $(minikube docker-env)

# Build the docker image
echo "Building Docker image..."
docker build -t golang-crud:v1 ..

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
kubectl apply -f configmap.yaml
kubectl apply -f persistence.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Wait for deployment to be ready
echo "Waiting for deployment to be ready..."
kubectl rollout status deployment/golang-crud

# Get the URL to access the service
echo "Getting service URL..."
minikube service golang-crud --url

echo "Deployment complete! Access your application at the URL above."
echo "For Swagger documentation, add '/swagger/index.html' to the URL."
