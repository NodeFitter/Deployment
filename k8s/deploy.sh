kubectl create configmap mariadb-initdb --from-file=initdb/

kubectl apply -f k8s/mariadb.yaml
kubectl apply -f k8s/app.yaml

# Allow/disallow scheduling on the control plane node for testing
kubectl taint nodes lorenzo-latitude-7280 node-role.kubernetes.io/control-plane:NoSchedule-
kubectl taint nodes lorenzo-latitude-7280 node-role.kubernetes.io/control-plane:NoSchedule
