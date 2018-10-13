kubectl get pods | grep $1 | awk '{print $1}' | xargs kubectl delete pods
