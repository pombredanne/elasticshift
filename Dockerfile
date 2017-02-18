FROM scratch
MAINTAINER Ghazni Nattarshah <ghazni.nattarshah@conspico.com>

COPY /bin/linux_386/armor /opt/conspico/armor/armor

EXPOSE 5050 5051

ENTRYPOINT ["/opt/conspico/armor/armor"]