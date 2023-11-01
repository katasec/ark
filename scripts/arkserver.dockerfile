# ------------------------------------------------------------------------------------------
# Build App in Golan Container
# ------------------------------------------------------------------------------------------

#FROM --platform=linux/amd64 golang:1.20.3-alpine as build
FROM golang:1.20.3-alpine as build
WORKDIR /go/src/app
COPY . .

RUN apk add gcc && \
    apk add musl-dev && \
    apk add libc-dev && \
    go mod download && \
    CGO_ENABLED=1 go build -o /go/bin/ark


# ------------------------------------------------------------------------------------------
# Copy compiled binary on to golang distro
# ------------------------------------------------------------------------------------------
#FROM --platform=linux/amd64 bash:5.2.15-alpine3.17
FROM bash:5.2.15-alpine3.17
COPY --from=build /go/bin/ark /usr/local/bin
#COPY --from=build /go/src/app/scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
CMD ["/ark","server"]


#docker run -it -v $HOME/.ark:/root/.ark ghcr.io/katasec/arkserver:v0.0.6