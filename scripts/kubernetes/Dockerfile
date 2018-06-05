FROM nginx:latest
MAINTAINER Ghazni Nattarshah <ghazni.nattarshah@conspico.com>

#COPY ./bin/linux_386/worker /opt/conspico

# update ca certs
#RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /etc/sudoers.d

# create elasticshift user/group
RUN groupadd --gid 1005 elasticshift \
	&& useradd --uid 1005 --gid elasticshift --shell /bin/bash --create-home elasticshift \
	&& echo 'elasticshift ALL=NOPASSWD: ALL' >> /etc/sudoers.d/elasticshift

RUN su - elasticshift

#EXPOSE 5050 5051

#ENTRYPOINT ["/opt/conspico/worker"]
