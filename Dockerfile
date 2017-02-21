FROM scratch
MAINTAINER Ghazni Nattarshah <ghazni.nattarshah@conspico.com>

COPY /bin/linux_386/elasticshift /opt/conspico/elasticshift

EXPOSE 5050 5051

ENTRYPOINT ["/opt/conspico/elasticshift/elasticshift"]