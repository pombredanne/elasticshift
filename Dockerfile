ARG ALPINE_VERSION=
FROM alpine:$ALPINE_VERSION

MAINTAINER Ghazni Nattarshah <ghazni.nattarshah@conspico.com>

COPY ./bin/linux_386/elasticshift /opt/conspico/

# update ca certs
RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["/opt/conspico/elasticshift"]
