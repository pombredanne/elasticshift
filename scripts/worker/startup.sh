#!/bin/bash

# set timezone to UTC
ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

if [ ! -d "/etc/sudoers.d" ]; then
    mkdir -p /etc/sudoers.d
fi

# create shiftuser user/group
groupadd --gid 1005 elasticshift \
    && useradd --uid 1005 --gid elasticshift --shell /bin/bash --create-home elasticshift \
    && echo elasticshift ALL=NOPASSWD: ALL >> /etc/sudoers.d/elasticshift

su elasticshift -c "${SHIFT_DIR}/sys/worker"