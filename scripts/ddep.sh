kubectl get deploy | grep $1 | awk '{print $1}' | xargs kubectl delete deploy
