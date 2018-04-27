FROM alpine:latest
MAINTAINER Ghazni Nattarshah <ghazni.nattarshah@conspico.com>

COPY ./bin/linux_386/elasticshift /opt/conspico/

# update ca certs
RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*

EXPOSE 5050 5051

ENTRYPOINT ["/opt/conspico/elasticshift"]
