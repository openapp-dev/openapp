#!/usr/bin/env bash

set -e

export KO_DOCKER_REPO=${KO_DOCKER_REPO:-localhost:5000}

source ./hack/utils.sh

# 1.Verify env
if [[ -z ${KO_DOCKER_REPO} ]]; then
    echo "KO_DOCKER_REPO is not set. Please set it to the docker repository where you want to store the image."
    exit 1
fi

# 2.Install go env
if [[ -z $(which go) ]]; then
    echo "go is not installed. Try to install go..."
    wget https://go.dev/dl/go1.20.6.linux-amd64.tar.gz -P /tmp/
    tar -C /usr/local -xzf /tmp/go1.20.6.linux-amd64.tar.gz
    rm -rf /tmp/go1.20.6.linux-amd64.tar.gz
fi
go mod tidy

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

# 5.Launch k3s and disable LoadBalancer
PROXY_ENV=""
if [[ -n ${https_proxy} ]]; then
    PROXY_ENV="-e HTTP_PROXY=${http_proxy} -e HTTPS_PROXY=${https_proxy} -e NO_PROXY=${no_proxy}"
fi
# check whether k3s container exist
if [[ -z $(docker ps -a --filter name=k3s -q) ]]; then
    docker run -itd --network host \
        --privileged \
        ${PROXY_ENV} \
        -v ${HOME}/openapp:/root/openapp \
        --name k3s rancher/k3s:v1.27.4-k3s1 server --disable traefik,servicelb

    # Give some time to let the k3s cluster initialize
    sleep 30

    mkdir -p ${HOME}/.config/
    docker cp k3s:/etc/rancher/k3s/k3s.yaml ${HOME}/.config/k3s.yaml
fi


# 6.Deploy openapp
export KUBECONFIG=${HOME}/.config/k3s.yaml
ko apply -Rf config

# 7.Echo related information
echo "Please run following command to control openapp:"
utils::echo_note "export KUBECONFIG=${HOME}/.config/k3s.yaml"
utils::echo_note "kubectl get apptemplate -A"
