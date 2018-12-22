#!/bin/sh

echo "SHIFT_DIR=${SHIFT_DIR}"
echo "WORKER_URL=${WORKER_URL}"

# set timezone to UTC
ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

if [ ! -d "/etc/sudoers.d" ]; then
    mkdir -p /etc/sudoers.d
fi

# create shiftuser user/group
if [ ! $(getent group elasticshift) ]; then
    groupadd --gid 1005 elasticshift \
        && useradd --uid 1005 --gid elasticshift --shell /bin/bash --create-home elasticshift \
        && echo elasticshift ALL=NOPASSWD: ALL >> /etc/sudoers.d/elasticshift
fi

if [ ! -f "${SHIFT_DIR}/sys/worker" ]; then
    wget -P "${SHIFT_DIR}/sys" -o worker "${WORKER_URL}/sys/worker" 2>/dev/null || curl -o "${SHIFT_DIR}/sys/worker" --create-dirs "${WORKER_URL}/sys/worker" 2>/dev/null
    chown -R elasticshift:elasticshift "${SHIFT_DIR}/sys/worker"
    chmod +x "${SHIFT_DIR}/sys/worker"
fi

su elasticshift -c "${SHIFT_DIR}/sys/worker"