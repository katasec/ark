# ------------------------------------------------------------------------------------------
# Start with Golang container
# - Build Ark binary
# - Download pulumi into container since ark for pulumi api automation
# ------------------------------------------------------------------------------------------
FROM golang:1.20.3-alpine as build
WORKDIR /go/src/app
COPY . .


RUN go mod download && \
    CGO_ENABLED=0 go build -o /go/bin/ark &&\
    apk add --no-cache curl && \
    curl -fsSL https://get.pulumi.com | sh && \
    /root/.pulumi/bin/pulumi plugin rm resource azure-native && \
    /root/.pulumi/bin/pulumi plugin install resource azure-native
    # wget https://get.pulumi.com/releases/sdk/pulumi-v3.90.1-linux-x64.tar.gz && \
    # tar -xzvf pulumi-v3.90.1-linux-x64.tar.gz --directory /tmp/ 

# ------------------------------------------------------------------------------------------
# Copy compiled binaries on to run container.
# ------------------------------------------------------------------------------------------

# FROM bash:5.2.15-alpine3.16
# FROM mcr.microsoft.com/dotnet/sdk:7.0

# Pulumi needs az cli
FROM mcr.microsoft.com/azure-cli:latest  
RUN apk add dotnet7-sdk
COPY --from=build /go/bin/ark /
#COPY --from=build /go/bin/ark /usr/local/bin
COPY --from=build /root/.pulumi/bin/pulumi /usr/local/bin/pulumi
COPY --from=build /root/.pulumi/bin/pulumi-language-dotnet /usr/local/bin/pulumi-language-dotnet
COPY --from=build /root/.pulumi/bin/pulumi-language-go /usr/local/bin/pulumi-language-go
CMD ["/ark"]
