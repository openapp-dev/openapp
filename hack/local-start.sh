#!/usr/bin/env bash

set -e

# 1.Verify env
if [[ -z ${KO_DOCKER_REPO} ]]; then
    echo "KO_DOCKER_REPO is not set. Please set it to the docker repository where you want to store the image."
    exit 1
fi

# 2.Install go env
if [[ -z $(which go) ]]; then
    echo "go is not installed. Try to install go..."
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -P /tmp/
    tar -C /usr/local -xzf /tmp/go1.22.0.linux-amd64.tar.gz
    rm -rf /tmp/go1.22.0.linux-amd64.tar.gz
fi

# 3.Install kubectl
if [[ -z $(which kubectl) ]]; then
    echo "kubectl is not installed. Try to install kubectl..."
    wget https://dl.k8s.io/release/v1.29.2/bin/linux/amd64/kubectl -P /tmp/
    chmod +x /tmp/kubectl
    mv /tmp/kubectl /usr/local/bin/kubectl
fi

# 4.Install ko
if [[ -z $(which ko) ]]; then
    echo "ko is not installed. Try to install ko..."
    wget https://github.com/ko-build/ko/releases/download/v0.15.2/ko_0.15.2_Linux_x86_64.tar.gz -P /tmp/
    tar -C /usr/local/bin/ -xzf /tmp/ko_0.15.2_Linux_x86_64.tar.gz
    rm -rf /tmp/ko_0.15.2_Linux_x86_64.tar.gz
fi

# 5.Install k3s
if [[ -z $(which k3s) ]]; then
    echo "k3s is not installed. Try to install k3s..."
    curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL=v1.28.7+k3s1 sh -
fi

# 6.Launch k3s and disable LoadBalancer
k3s-killall.sh
nohup k3s server --disable traefik,servicelb > /var/log/k3s.log &
sleep 30

# 7. Deploy openapp
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
ko apply -Rf config

# 8. Echo related information
echo "Please run following command to control openapp:"
echo "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml"
echo "kubectl get apptemplate -A"