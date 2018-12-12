#!/usr/bin/env bash

echo "SHIFT_DIR=${SHIFT_DIR}"
echo "WORKER_URL=${WORKER_URL}"

init() {

  # set the arch
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="armv7";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac

  # set the os
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
  
  # set timezone to UTC
  ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime

  #set sys dir
  SHIFT_SYS_DIR="${SHIFT_DIR}/sys"
}

createUserAndGroup() {

    if [ ! -d "/etc/sudoers.d" ]; then
        mkdir -p /etc/sudoers.d
    fi

    # create shiftuser user/group
    if [ ! $(getent group elasticshift) ]; then
        groupadd --gid 1005 elasticshift \
            && useradd --uid 1005 --gid elasticshift --shell /bin/bash --create-home elasticshift \
            && echo elasticshift ALL=NOPASSWD: ALL >> /etc/sudoers.d/elasticshift
    fi
}

downloadWorker() {
    WORKER_DIST="worker-$OS-$ARCH.tar.gz"
    if [ ! -f "${SHIFT_SYS_DIR}/worker" ]; then
        if type "curl" > /dev/null; then
            curl -SsL --create-dirs "${WORKER_URL}" -o "${SHIFT_SYS_DIR}/${WORKER_DIST}"
        elif type "wget" > /dev/null; then
            wget -P "${SHIFT_SYS_DIR}" -q -O "$WORKER_DIST" "${WORKER_URL}"
        fi

        chown -R elasticshift:elasticshift "${SHIFT_SYS_DIR}"

        tar xf "${SHIFT_SYS_DIR}/${WORKER_DIST}" -C "${SHIFT_SYS_DIR}"
        
        chmod +x "${SHIFT_SYS_DIR}/worker"
    fi
}

startWorker() {
    su elasticshift -c "${SHIFT_SYS_DIR}/worker"
}

#Stop execution on any error
#trap "fail_trap" EXIT
set -e

init
createUserAndGroup
downloadWorker
startWorker