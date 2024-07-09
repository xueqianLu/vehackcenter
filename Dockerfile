FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

WORKDIR /build

ADD . /build/vehackcenter

RUN  cd /build/vehackcenter && make && cp build/bin/vecenter /vecenter

FROM alpine

WORKDIR /root

COPY  --from=build /vecenter /usr/bin/vecenter

EXPOSE 9000

ENTRYPOINT [ "vecenter" ]
