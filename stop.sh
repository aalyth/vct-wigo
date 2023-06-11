kubectl delete service wigo
kubectl delete hpa wigo
kubectl delete deployment wigo --cascade=foreground
docker image rm -f wigo-img
