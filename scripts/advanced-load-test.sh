#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to monitor resources
monitor_resources() {
    local duration=$1
    local interval=5
    local end=$((SECONDS + duration))
    
    echo -e "${YELLOW}Starting resource monitoring for $duration seconds${NC}"
    
    printf "%-20s %-15s %-15s %-15s %-15s\n" "TIMESTAMP" "CPU(cores)" "MEMORY(Mi)" "PODS" "RESTARTS"
    
    while [ $SECONDS -lt $end ]; do
        local cpu=$(kubectl top pods -l app=service1 --no-headers | awk '{sum+=$2}END{print sum}')
        local memory=$(kubectl top pods -l app=service1 --no-headers | awk '{sum+=$3}END{print sum}')
        local pods=$(kubectl get pods -l app=service1 --no-headers | wc -l)
        local restarts=$(kubectl get pods -l app=service1 --no-headers | awk '{sum+=$4}END{print sum}')
        
        printf "%-20s %-15s %-15s %-15s %-15s\n" "$(date +%T)" "$cpu" "$memory" "$pods" "$restarts"
        sleep $interval
    done
}

# Function for detailed latency testing
test_latency() {
    local requests=$1
    local concurrent=$2
    
    echo -e "${GREEN}Running latency test with $requests requests ($concurrent concurrent connections)${NC}"
    
    ab -n $requests -c $concurrent -H "Host: app1.test" "http://$(minikube ip)/app1/" > ab_results.txt
    
    # Extract and format results
    echo -e "\n${YELLOW}Latency Results:${NC}"
    echo "----------------------------------------"
    grep "Requests per second" ab_results.txt
    grep "Time per request" ab_results.txt
    grep "Failed requests" ab_results.txt
    
    echo -e "\n${YELLOW}Latency Distribution:${NC}"
    echo "----------------------------------------"
    sed -n '/Percentage of the requests served/,/100%/p' ab_results.txt
}

# Function to check traffic distribution
check_traffic_distribution() {
    local samples=$1
    local v1=0
    local v2=0
    
    echo -e "${GREEN}Checking traffic distribution with $samples samples${NC}"
    
    for i in $(seq 1 $samples); do
        response=$(curl -s -H "Host: app1.test" "http://$(minikube ip)/app1")
        if echo "$response" | grep -q "Version 2"; then
            ((v2++))
        else
            ((v1++))
        fi
        echo -n "."
        if [ $((i % 50)) -eq 0 ]; then
            echo " $i"
        fi
    done
    
    local v1_percent=$((v1 * 100 / samples))
    local v2_percent=$((v2 * 100 / samples))
    
    echo -e "\n\n${YELLOW}Final Distribution:${NC}"
    echo "----------------------------------------"
    echo "Version 1: $v1 requests ($v1_percent%)"
    echo "Version 2: $v2 requests ($v2_percent%)"
}

# Main
echo -e "${GREEN}Starting comprehensive load and distribution test${NC}"

# 1. Monitor resources in background
monitor_resources 60 &
monitor_pid=$!

# 2. Run latency test
test_latency 1000 50

# 3. Check traffic distribution
check_traffic_distribution 200

# Wait for monitoring to finish
wait $monitor_pid

# Show final pod status
echo -e "\n${YELLOW}Final pod status:${NC}"
kubectl get pods -l app=service1 -o wide

echo -e "\n${GREEN}Test complete!${NC}"