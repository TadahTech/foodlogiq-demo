FROM alpine:latest as alpine
RUN apk add --no-cache tzdata ca-certificates

FROM golang:1.18.2-alpine as gobuild

ENV CGO_ENABLED 0

WORKDIR /opt

# COPY mod file and download dependencies so
# it can be cached as a layer to speed up
# successive builds
COPY go.mod .
COPY go.sum .
RUN go mod download

# COPY the source code as the last step
COPY . .

ADD . .
RUN go build -o build/foodlogiq-demo -ldflags "-X main.version=0.0.1-DEV" ./cmd/foodlogiq/

FROM scratch as release

COPY --from=alpine /etc/ssl/certs /etc/ssl/certs
COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=gobuild /opt/build/foodlogiq-demo /opt/foodlogiq-demo

WORKDIR /opt
EXPOSE 8000

ENTRYPOINT ["/opt/foodlogiq-demo"]

