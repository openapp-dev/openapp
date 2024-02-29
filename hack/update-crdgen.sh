#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

controller-gen crd paths=./pkg/apis/app/... output:crd:dir=./config/crds
controller-gen crd paths=./pkg/apis/service/... output:crd:dir=./config/crds
