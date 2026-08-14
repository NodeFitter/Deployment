# Deployment
Repository containing files for the presentation of the project

## How to setup custom metrics for HPA
- Install `Helm` with:
    ```sh
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-4
    chmod +x get_helm.sh
    ./get_helm.sh
    ```
- Add the [Prometheus Community Kubernetes Helm Charts](https://github.com/prometheus-community/helm-charts) with:
    ```sh
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    ```
- Install the `kube-prometheus-stack` without unnecessary component and set the namespace to `monitoring`
    ```sh
    helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace --set grafana.enabled=false --set alertmanager.enabled=false --set kubeStateMetrics.enabled=false --set nodeExporter.enabled=false --set prometheus-pushgateway.enabled=false --set defaultRules.create=false --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false --set prometheus.prometheusSpec.serviceMonitorSelector.matchLabels.monitoring=prometheus
    ```
- Install the `prometheus-adapter`, necessary to make possible the retrieval of measurements by kubernetes:
    ```sh
    helm install prometheus-adapter prometheus-community/prometheus-adapter --namespace monitoring --set prometheus.url=http://prometheus-kube-prometheus-prometheus.monitoring.svc --set prometheus.port=9090 -f ./k8s/prometheusConfig.yaml 
    ```
