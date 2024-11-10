#!/bin/bash

# To send requests
send_requests() {
    local host=$1
    local path=$2
    local requests=$3
    
    echo "Sending $requests requests to $host$path"
    for i in $(seq 1 $requests); do
        curl -s -H "Host: $host" "http://$MINIKUBE_IP$path" > /dev/null &
        if [ $((i % 10)) -eq 0 ]; then
            echo -n "."
        fi
    done
    echo "Finished!"
}

# Configure variables
MINIKUBE_IP=$(minikube ip)

# Basic load test
echo "Initializing basic load test..."
send_requests app1.test /app1 100
send_requests app2.test /app2 100

# Latency test
echo -e "\nTesting latency..."
time curl -H "Host: app1.test" http://$MINIKUBE_IP/app1
time curl -H "Host: app2.test" http://$MINIKUBE_IP/app2

# Concurrency test
echo -e "\nTesting concurrency..."
ab -n 1000 -c 50 -H "Host: app1.test" http://$MINIKUBE_IP/app1/