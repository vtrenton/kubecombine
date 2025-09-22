# Kubecombine
Combine several kubeconfigs together in one file!

## Building
There are three main ways to build this:
1. good old go build
2. The nix way (including the nix docker way)
3. The Docker way

### Go build:
`go build -o kubecombine cmd/kubecombine/combine.go`

### Nix
build a binary directly with:
`nix build`
The binary will then exist at:
`result/bin/kubecombine`

OR you can build a docker image and load it with docker
Protip: the nix image is going to be smaller than the regular docker one :)
```
nix build .#docker
docker load < result
```
### Docker
Building can be accomplished with:
`docker build`


## Useage

### Docker

`  docker run -v /path/to/configs:/configs kubecombine:latest /configs/config1.yaml /configs/config2.yaml`
