#!/bin/bash

# Function to wait for pods
wait_for_pods() {
    echo "Waiting for pods to become ready..."
    kubectl wait --for=condition=ready pod -l app=service1 --timeout=120s
}

# Clean up existing resources
echo "Cleaning up previous resources..."
kubectl delete -f custom-nginx.yaml --ignore-not-found
kubectl delete -f canary.yaml --ignore-not-found

# Apply base configurations
echo "Applying base configurations..."
kubectl apply -f custom-nginx.yaml
wait_for_pods

echo "Applying canary configurations..."
kubectl apply -f canary.yaml
wait_for_pods

# Function to test distribution
test_distribution() {
    local total_requests=100
    local v1_count=0
    local v2_count=0
    
    echo "Testing traffic distribution with $total_requests requests..."
    
    for i in $(seq 1 $total_requests); do
        response=$(curl -s -H "Host: app1.test" "http://$(minikube ip)/app1")
        if echo "$response" | grep -q "Version 2"; then
            ((v2_count++))
            echo -n "v2."
        else
            ((v1_count++))
            echo -n "v1."
        fi
        
        if [ $((i % 10)) -eq 0 ]; then
            echo " "
        fi
    done
    
    echo -e "\nResults:"
    echo "Version 1: $v1_count requests ($(( v1_count * 100 / total_requests ))%)"
    echo "Version 2: $v2_count requests ($(( v2_count * 100 / total_requests ))%)"
}

# Display status
echo -e "\nPod status:"
kubectl get pods -l app=service1 -o wide

echo -e "\nService status:"
kubectl get services

echo -e "\nIngress status:"
kubectl get ingress

# Run distribution test
echo -e "\nStarting distribution test..."
test_distribution

# Load test with ab (if installed)
if command -v ab &> /dev/null; then
    echo -e "\nRunning load test..."
    ab -n 1000 -c 50 -H "Host: app1.test" "http://$(minikube ip)/app1/"
else
    echo -e "\nApache Benchmark (ab) not found. Install: sudo apt-get install apache2-utils"
fi