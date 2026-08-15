#!/usr/bin/env bash

set -e

K8S_DIR="k8s"

INITDB_DIR="initdb"
INITDB_CONFIGMAP_NAME="mariadb-initdb"

DB_SECRET_FILE="secrets/mariadb.env"
DB_SECRET_NAME="mariadb-secret"
APP_SECRET_FILE="secrets/app.env"
APP_SECRET_NAME="app-secret"

cleanup() {
    echo
    echo "Cleaning up..."

    echo "Deleting application..."
    kubectl delete deployment app --ignore-not-found --wait=false
    kubectl delete service app --ignore-not-found --wait=false

    echo "Deleting MariaDB..."
    kubectl delete statefulset mariadb --ignore-not-found --wait=false
    kubectl delete service mariadb --ignore-not-found --wait=false

    echo "Deleting MariaDB PVC..."
    kubectl delete pvc data-mariadb-0 --ignore-not-found --wait=false

    echo "Deleting MariaDB PV..."
    kubectl delete pv mariadb-pv --ignore-not-found --wait=false

    echo "Deleting configMap..."
    kubectl delete configmap "$INITDB_CONFIGMAP_NAME" --ignore-not-found

    echo "Deleting secrets..."
    kubectl delete secret "$DB_SECRET_NAME" --ignore-not-found
    kubectl delete secret "$APP_SECRET_NAME" --ignore-not-found

    echo "Deleting namespaces..."
    kubectl delete namespace frontend
    kubectl delete namespace backend

    echo "Cleanup complete."
}

trap cleanup EXIT

echo "Creating namespaces..."
kubectl create namespace frontend
kubectl create namespace backend

echo "Creating the config map..."
kubectl create configmap "$INITDB_CONFIGMAP_NAME" \
    --from-file="$INITDB_DIR" \
    --dry-run=client \
    -o yaml |
    kubectl apply -f -

echo "Creating secrets..."

kubectl create secret generic "$DB_SECRET_NAME" \
    --from-env-file="$DB_SECRET_FILE" \
    --dry-run=client \
    -o yaml |
    kubectl apply -f -

kubectl create secret generic "$APP_SECRET_NAME" \
    --from-env-file="$APP_SECRET_FILE" \
    --dry-run=client \
    -o yaml |
    kubectl apply -f -


echo "Applying manifests..."

kubectl apply -f "$K8S_DIR/"

echo
echo "Deployment is running."
echo "Press ENTER to delete everything and clean up."
read -r

cleanup
