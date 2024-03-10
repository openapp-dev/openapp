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

Or if you don't want use docker hub, you can use the following command:
```sh
git clone https://github.com/openapp-dev/openapp.git
cd openapp

docker run -d -p 5000:5000 --name registry registry:2.7
export KO_DOCKER_REPO=localhost:5000
hack/local-start.sh
```

## Testing OpenAPP API

Run following command to test the basic API:
```
curl --silent http://localhost:30003/api/v1/apps/templates | jq .
```

## Deploy with code updates

Once you modified the code, just run the following command to re-deploy the modified code:
```sh
hack/local-start.sh
```

## Give a PR

Once you finish developing and do enough testing, you can create a PR in this repository. The maintainer will check and merge it as soon as possible.