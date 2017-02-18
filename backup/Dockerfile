FROM golang:1.7-alpine

ADD . /go/src/gitlab.com/conspico/esh

RUN go build /go/src/gitlab.com/conspico/esh/cmd/esh/main.go

#ENTRYPOINT /go/bin/esh

EXPOSE 80
EXPOSE 443