FROM golang AS builder
WORKDIR /kubecombine
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o kubecombine cmd/kubecombine/combine.go

FROM alpine
COPY --from=builder /kubecombine/kubecombine /usr/local/bin/kubecombine
ENTRYPOINT ["/usr/local/bin/kubecombine"]
