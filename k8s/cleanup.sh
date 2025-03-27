#!/bin/bash

echo "Cleaning up GoLang-CRUD Kubernetes resources..."

# Delete Kubernetes resources
kubectl delete service golang-crud || true
kubectl delete deployment golang-crud || true
kubectl delete pvc golang-crud-pvc || true
kubectl delete pv golang-crud-pv || true
kubectl delete configmap golang-crud-config || true

echo "Cleanup complete!"
