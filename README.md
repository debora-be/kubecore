# kubecore

a cli tool for managing and testing kubernetes microservices locally using minikube

## overview

launch, manage, and monitor a small ecosystem running on a cluster. it provides functionality for deployment, scaling, and status checking of services, along with features for testing service mesh configurations like traffic routing and load balancing.

## features

- **service deployment**: easily deploy microservices with `kubecore deploy service-name`
- **scaling**: adjust the number of replicas with `kubecore scale service-name --replicas=3`
- **status monitoring**: check service status, pods, and metrics with `kubecore status`
- **ingress management**: create and manage ingress rules for your services
- **canary deployments**: test new versions of your services with traffic splitting

## prerequisites

- go 1.21 or higher
- minikube
- kubectl
- docker
- apache benchmark (ab) for load testing

## installation

1. build the project:
```bash
go build -o kubecore
```

## usage

### starting the environment

1. start minikube:
```bash
minikube start --cpus=2 --memory=4096mb --driver=docker
```

2. enable ingress:
```bash
minikube addons enable ingress
```

### basic commands

1. deploy a service:
```bash
./kubecore deploy myservice
```

2. scale the service:
```bash
./kubecore scale myservice --replicas=3
```

3. check status:
```bash
./kubecore status
```

### testing scenarios

#### 1. basic service deployment
```bash
# deploy two services
./kubecore deploy service1
./kubecore deploy service2

# verify deployment
./kubecore status
```

#### 2. ingress testing
```bash
# create ingress rules
kubectl apply -f - <<EOF
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

# add to /etc/hosts
echo "$(minikube_ip) app1.test" | sudo tee -a /etc/hosts

# test the endpoint
curl -H "Host: app1.test" http://app1.test/app1
```

#### 3. canary deployment
```bash
# deploy v1
./kubecore deploy service1-v1

# deploy v2 with traffic splitting
kubectl apply -f canary.yaml

# test traffic distribution
./advanced-load-test.sh
```

#### 4. load testing
```bash
# install apache benchmark if not installed
sudo apt-get install apache2-utils

# run load test
ab -n 1000 -c 50 -H "Host: app1.test" http://(MINIKUBE_IP)/app1/
```

## monitoring and debugging

### view logs
```bash
# ingress controller logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller -f

# service logs
kubectl logs -l app=service1
```

### check metrics
```bash
# enable metrics server
minikube addons enable metrics-server

# view pod metrics
kubectl top pods
```

## running tests

1. traffic test:
```bash
./scripts/traffic-test.sh
```

2. canary deployment test:
```bash
./scripts/canary-test.sh
```

3. advanced load test:
```bash
./scripts/advanced-load-test.sh
```

4. ingress test:
```bash
./scripts/ingress-test.sh
```

## troubleshooting

### common issues

1. ingress not working:
```bash
# verify ingress controller is running
kubectl get pods -n ingress-nginx

# check ingress configuration
kubectl describe ingress service1-ingress
```

2. service not responding:
```bash
# check pod status
kubectl get pods -l app=service1

# check service endpoints
kubectl describe service service1
```

3. performance issues:
```bash
# check resource usage
kubectl top pods
kubectl top nodes
```