# ------------------------------------------------------------------------------------------
# Start with Golang container
# - Build Ark binary
# - Download pulumi into container since ark for pulumi api automation
# ------------------------------------------------------------------------------------------
FROM golang:1.18.4-alpine as build
WORKDIR /go/src/app
COPY . .

RUN go mod download && \
    CGO_ENABLED=0 go build -o /go/bin/ark &&\
    wget https://get.pulumi.com/releases/sdk/pulumi-v3.52.0-linux-x64.tar.gz && \
    tar -xzvf pulumi-v3.52.0-linux-x64.tar.gz --directory /tmp/ 

# ------------------------------------------------------------------------------------------
# Copy compiled binaries on to run container.
# ------------------------------------------------------------------------------------------

# FROM bash:5.2.15-alpine3.16

FROM mcr.microsoft.com/azure-cli:latest
COPY --from=build /go/bin/ark /
COPY --from=build /go/bin/ark /usr/local/bin
COPY --from=build /tmp/pulumi/pulumi /usr/local/bin/pulumi
CMD ["/ark"]