#!/usr/bin/env bash

CLUSTER_NAME="test-cluster"
REGISTRY_NAME="test-registry"

# Get the current Kubernetes context to switch back once done
PREV_KUBE_CONTEXT=$(kubectl config current-context)

# Check if test cluster already exists
MATCHED_CLUSTERS=$(k3d cluster list | grep $CLUSTER_NAME | wc -l)
if [ $MATCHED_CLUSTERS -gt 0 ]; then
    echo "error: $CLUSTER_NAME already exists"
    exit 1
fi

# Check if test registry already exists
MATCHED_REGISTRIES=$(k3d registry list | grep $REGISTRY_NAME | wc -l)
if [ $MATCHED_REGISTRIES -gt 0 ]; then
  echo "error: $REGISTRY_NAME already exists"
  exit 1
fi

echo "creating registry $REGISTRY_NAME"
k3d registry create $REGISTRY_NAME.localhost --port 5000

echo "creating cluster $CLUSTER_NAME"
k3d cluster create $CLUSTER_NAME --config yaml/cluster-config.yaml --registry-use k3d-$REGISTRY_NAME.localhost:5000 --registry-config yaml/registries.yaml

# Use test cluster context
kubectl config use-context k3d-$CLUSTER_NAME

echo "building Docker images"
docker build -t localhost:5000/web-client/api ../../
docker build -t localhost:5000/web-client/frontend ../../../frontend

echo "pushing Docker images"
docker push localhost:5000/web-client/api
docker push localhost:5000/web-client/frontend

# Install HAProxy
helm repo add haproxytech https://haproxytech.github.io/helm-charts
helm repo update
helm install haproxy haproxytech/kubernetes-ingress \
  --set-string "controller.extraArgs={--sync-period=60s}" \
  --set controller.startupProbe.initialDelaySeconds=80

echo "creating secrets frontend-secret and backend-secret"
kubectl create secret generic frontend-secret --from-env-file=.frontend.env
kubectl create secret generic backend-secret --from-env-file=.backend.env

# Apply backend environment variables
set -o allexport && source .backend.env && set +o allexport

echo "creating PostgreSQL database"
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-postgresql bitnami/postgresql --version 16.0.6 --set global.postgresql.auth.password=$DB_PASSWORD

echo "creating API Deployment"
kubectl create -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels:
    k8s-app: api
spec:
  selector:
    matchLabels:
      k8s-app: api
  template:
    metadata:
      labels:
        k8s-app: api
    spec:
      containers:
      - name: api
        image: localhost:5000/web-client/api
        resources:
          limits:
            memory: 1Gi
            cpu: 1
        ports:
        - containerPort: 8090
        env:
        - name: ENV
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: ENV
        - name: PORT
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: PORT
        - name: OAUTH_CLIENT_ID
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: OAUTH_CLIENT_ID
        - name: OAUTH_CLIENT_SECRET
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: OAUTH_CLIENT_SECRET
        - name: OAUTH_CALLBACK_URL
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: OAUTH_CALLBACK_URL
        - name: ISSUER
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: ISSUER
        - name: API_URL
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: API_URL
        - name: DOMAIN
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: DOMAIN
        - name: APP_URL
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: APP_URL
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: DB_URL
        - name: CLUSTER_TYPE
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: CLUSTER_TYPE
        - name: KUBE_TOKEN
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: KUBE_TOKEN
        - name: KUBE_CERT
          valueFrom:
            secretKeyRef:
              name: backend-secret
              key: KUBE_CERT
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
EOF

kubectl get all

# Wait for HAProxy deployment to be ready
kubectl wait deployment haproxy-kubernetes-ingress --for condition=Available=True --timeout=180s

echo "creating API service"
kubectl apply -f yaml/api-svc.yaml

echo "creating ingress"
kubectl apply -f yaml/ingress.yaml

PORT_FORWARD=8082
echo "establishing port-forward on port ${PORT_FORWARD}"
kubectl port-forward service/haproxy-kubernetes-ingress $PORT_FORWARD:80 &

read -n 1 -s

echo "stopping port-forward"
pkill kubectl

echo "deleting resources"
kubectl delete -f yaml/ingress.yaml
kubectl delete -f yaml/api-svc.yaml
kubectl delete deployment api

# Clean up
k3d cluster delete $CLUSTER_NAME
k3d registry delete k3d-$REGISTRY_NAME.localhost

# Switch back to original Kubernetes context
kubectl config use-context $PREV_KUBE_CONTEXT