docker image rm -f wigo-img
docker build -t wigo-img .

kubectl delete -f components.yaml
kubectl apply -f k8s/components.yaml # this is so there're no problems with the metrics-server

kubectl delete deployment wigo --cascade=foreground
kubectl apply -f k8s/wigo.yaml

kubectl delete hpa wigo
kubectl autoscale deployment wigo --cpu-percent=60 --min=3 --max=15

kubectl delete service wigo-svc
kubectl expose deployment wigo --name=wigo-svc --type=LoadBalancer --port=80 --target-port=4000 

kubectl delete ing wigo-ing
# kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.0/deploy/static/provider/cloud/deploy.yaml
# kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s
# kubectl apply -f k8s/ingress.yaml

kubectl get hpa wigo -w
