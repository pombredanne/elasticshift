#!/bin/bash

export PATH=$PATH:$HOME/gopath/bin

OPTIONS="-config=config.yml -env develop"

set -ex

sql-migrate status $OPTIONS
sql-migrate up $OPTIONS
sql-migrate status $OPTIONS