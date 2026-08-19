# Deployment
Repository containing files for the presentation of the project

> [!NOTE]
> If the various helm commands are run where the Kubernetes control plane is, installation may fail due to the control plane being tainted.
> It is possible to un-taint the node by running
> ```sh
> kubectl taint nodes <node-name> node-role.kubernetes.io/control-plane:NoSchedule-
> ```

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
    helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace -f ./helmCharts/prometheusChart.yaml
    ```
- Install the `prometheus-adapter`, necessary to make possible the retrieval of measurements by kubernetes:
    ```sh
    helm install prometheus-adapter prometheus-community/prometheus-adapter --namespace monitoring -f ./helmCharts/prometheusAdapterChart.yaml 
    ```
- Add a local path provisioner (in this example [rancher](https://github.com/rancher/local-path-provisioner) is used):
    ```sh
    kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
    ```
- Install `loki` as the log collector
    ```sh
    helm repo add grafana-community https://grafana-community.github.io/helm-charts
    helm repo update
    helm install loki grafana-community/loki -f ./helmCharts/lokiChart.yaml -n monitoring
    ```
- Install `fluent-bit` as the log collector (daemonset)
    ```sh
        helm repo add fluent https://fluent.github.io/helm-charts
        helm upgrade --install fluent-bit fluent/fluent-bit -n monitoring -f ./helmCharts/fluentBitChart.yaml
    ```