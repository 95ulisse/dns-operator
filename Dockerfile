# Build the operator binary
FROM golang:1.13-alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o dns-operator main.go

# Final image is based on Alpine
FROM alpine:3
WORKDIR /
COPY --from=builder /workspace/dns-operator .
USER nobody:nobody

ENTRYPOINT ["/dns-operator"]
