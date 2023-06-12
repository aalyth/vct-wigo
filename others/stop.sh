kubectl delete service wigo-svc
kubectl delete hpa wigo
kubectl delete deployment wigo --cascade=foreground
docker image rm -f wigo-img
