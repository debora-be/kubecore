#!/bin/bash

# Function to wait for the service to be ready
wait_for_service() {
    echo "Standing by service $1 to be ready..."
    kubectl wait --for=condition=available --timeout=60s deployment/$1
}

# Check if Minikube is running
if ! minikube status > /dev/null 2>&1; then
    echo "Initializing Minikube..."
    minikube start --cpus=2 --memory=4096mb --driver=docker
fi

# Tur on the Ingress addon
echo "Turning on Ingress Controller..."
minikube addons enable ingress

# Deploy test services
echo "Deploying test services..."
./microkube deploy service1
wait_for_service service1

./microkube deploy service2
wait_for_service service2

# Create Ingress for service1
echo "Creating Ingress for service1..."
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: service1-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: app1.test
    http:
      paths:
      - path: /app1
        pathType: Prefix
        backend:
          service:
            name: service1
            port:
              number: 80
EOF

# Create Ingress for service2
echo "Creating Ingress for service2..."
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: service2-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: app2.test
    http:
      paths:
      - path: /app2
        pathType: Prefix
        backend:
          service:
            name: service2
            port:
              number: 80
EOF

# Stand by Ingress Controller to be ready
echo "Standing by Ingress Controller..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# Get Minikube IP
MINIKUBE_IP=$(minikube ip)
echo "Minikube IP: $MINIKUBE_IP"

# Commands for /etc/hosts
echo "Run the following command to add entries at /etc/hosts:"
echo "sudo sh -c \"echo '$MINIKUBE_IP app1.test app2.test' >> /etc/hosts\""

# Wait for a while to Ingress propagate
echo "Standing by Ingress propagation..."
sleep 10

# Show Ingress status
echo "Ingress status:"
kubectl get ingress

# Test endpoints
echo "To test endpoints, run:"
echo "curl -H \"Host: app1.test\" http://$MINIKUBE_IP/app1"
echo "curl -H \"Host: app2.test\" http://$MINIKUBE_IP/app2"

echo "Setup complete!"