# Kubecombine
Combine several kubeconfigs together in one file!

## Building
There are three main ways to build this:
1. good old go build
2. The nix way (including the nix docker way)
3. The Docker way

### Go:
`go build -o kubecombine cmd/kubecombine/combine.go`

### Nix
build a binary directly with:
`nix build`
The binary will then exist at:
`result/bin/kubecombine`

OR you can build a docker image and load it with docker
```
nix build .#docker
docker load < result
```
### Docker
Building can be accomplished with:
`docker build . -t kubecombine`


## Usage

### Docker

`  docker run -v /path/to/configs:/configs kubecombine:latest /configs/config1.yaml /configs/config2.yaml`
