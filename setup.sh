docker image rm -f wigo-img
docker build -t wigo-img .

kubectl delete -f components.yaml
kubectl apply -f k8s/components.yaml # this is so there're no problems with the metrics-server

kubectl delete deployment wigo
kubectl apply -f k8s/wigo.yaml

kubectl delete hpa wigo
kubectl autoscale deployment wigo --cpu-percent=60 --min=3 --max=12

kubectl delete service wigo
kubectl expose deployment wigo --type=LoadBalancer --port=8080 --target-port=4000
