# Contribution guidelines

If you want to contribute to OpenAPP, please read the following content.

## Install the development environment

Before developing, you need to init the environment.
```sh
git clone https://github.com/openapp-dev/openapp.git
cd openapp

export KO_DOCKER_REPO=<your dockerhub repo>
hack/local-start.sh
```

## Deploy with code updates

Once you modified the code, just run the following command to re-deploy the modified code:
```sh
hack/local-start.sh
```

## Give a PR

Once you finish developing and do enough testing, you can create a PR in this repository. The maintainer will check and merge it as soon as possible.