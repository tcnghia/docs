A simple gRPC server written in Go that you can use for testing.

## Prerequisites

- [Install the latest version of Knative Serving](../../../install/README.md).

- Install [docker](https://www.docker.com/).

- Download a copy of the code:

  ```shell
  git clone -b "release-0.7" https://github.com/knative/docs knative-docs
  cd knative-docs/docs/serving/samples/grpc-ping-go
  ```

## Build and run the gRPC server

First, build and publish the gRPC server to DockerHub:

```shell
# Build and publish the container, run from the root directory.
docker build \
  --tag "docker.io/tcnghia/grpc-ping-go" \
  --file=docs/serving/samples/grpc-ping-go/Dockerfile .
docker push "docker.io/tcnghia/grpc-ping-go"
```


```shell
kubectl apply --filename docs/serving/samples/grpc-ping-go/sample.yaml
```

## Use the client to stream messages to the gRPC server

1. Fetch the Service's hostname (assumes DNS has been configured).

   ```shell
   # Put the ingress IP into an environment variable.
   export GRPC_URL=$(basename $(kubectl get ksvc grpc-knative --output jsonpath="{.status.url}"))
   ```

1. Use the client to send message streams to the gRPC server

   ```shell
   docker run -ti --entrypoint=/client docker.io/tcnghia/grpc-ping-go \
     -server="${GRPC_URL}"
```
