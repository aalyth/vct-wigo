docker build -t wigo-img .

minikube start
timeout 30 > NUL
minikube image load wigo-img
kubectl delete deployment wigo
kubectl apply -f k8s/wigo.yaml
kubectl delete service wigo
kubectl autoscale deployment wigo --cpu-percent=60 --min=3 --max=12
kubectl expose deployment wigo --type=LoadBalancer --port=4000
timeout 15 > NUL
kubectl port-forward svc/wigo 80:4000
