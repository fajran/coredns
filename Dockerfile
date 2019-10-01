FROM golang:1.12 AS build

WORKDIR /go/src/github.com/coredns/coredns

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install github.com/coredns/coredns

FROM debian:stable-slim AS ca

RUN apt-get update && apt-get -uy upgrade
RUN apt-get -y install ca-certificates && update-ca-certificates

FROM scratch

COPY --from=ca /etc/ssl/certs /etc/ssl/certs
COPY --from=build /go/bin/coredns /coredns

EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]

