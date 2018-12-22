FROM alpine:latest
COPY ./startup.sh /tmp/shiftsys
RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT [ "/tmp/shiftsys/startup.sh" ]
