#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "$0" )" && pwd )"

PID_FILE="$DIR/servers.pid"

if [ -f "$PID_FILE" ]; then
    echo "servers.pid exists!"
    exit 1
fi

for port in 8008 8009; do
    rm -rf "${DIR:?}/$port"
    rm -rf "$DIR/media_store.$port"
done

rm -rf "${DIR:?}/etc"

