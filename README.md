# Deployment
Repository containing files for the presentation of the project

> [!NOTE]
> This part of the setup, not necessary to use NodeFitter but useful to replicate what demonstrated during the project presentation, uses Helm to install some components. To install Helm, run the following commands:
> ```sh
> curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-4
> chmod +x get_helm.sh
> ./get_helm.sh
> ```

> [!IMPORTANT]
> If the various Helm commands are run where the Kubernetes control plane is, installation may fail due to the control plane being tainted.
> It is possible to un-taint the node by running
> ```sh
> kubectl taint nodes <node-name> node-role.kubernetes.io/control-plane:NoSchedule-
> ```

> [!NOTE]
> Contrary to the installation of kubernetes with `minikube`, the standard installation of `kubeadm` does not include any dynamic provisioner for PersistentVolumes, therefore, PersistentVolumes should be manually configured or it is possible to install a local dynamic provisioner like [rancher local-path-provisioner](https://github.com/rancher/local-path-provisioner).
> To install the mentioned provisioner via Helm, run:
> ```sh
> kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
> ```
> The Helm values under the [helmCharts](https://github.com/NodeFitter/Deployment/tree/main/helmCharts) consider the mentioned provider as already installed, as well as the PersistentVolumeClaim of Grafana, see the [monitoring.yml](https://github.com/NodeFitter/Deployment/blob/main/k8s/monitoring.yml) file


## Setup custom metrics for HPA

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

## Setup logs collection

- Add the [Grafana Community Helm Charts](https://grafana-community.github.io/helm-charts/):
  ```sh
  helm repo add grafana-community https://grafana-community.github.io/helm-charts
  helm repo update
  ```

- Install `loki` to centralize logs:
    ```sh
    helm install loki grafana-community/loki -f ./helmCharts/lokiChart.yaml -n monitoring
    ```
- Install `fluent-bit` (daemonset) to collect logs from the various nodes:
    ```sh
    helm repo add fluent https://fluent.github.io/helm-charts
    helm upgrade --install fluent-bit fluent/fluent-bit -n monitoring -f ./helmCharts/fluentBitChart.yaml
    ```