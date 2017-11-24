#!/usr/bin/env bash

set -e

ROOT="$( cd $( dirname $0 ) && pwd -P )/.."
VENDOR="$ROOT/vendor"

GOPATH0="$ROOT/.gopath"
GOPATH="$VENDOR:$GOPATH0:$GOPATH"
GOPATH="$( echo $GOPATH | sed -e 's/:$//g' )" # drop tralling comma
echo "GOPATH: $GOPATH"
export GOPATH

export TMPDIR="$ROOT/.tmp"
export CGO_ENABLED="0" # 有效減少 dependencies
export BuildTags="$BuildTags exclude_graphdriver_devicemapper exclude_graphdriver_aufs exclude_graphdriver_btrfs" # 有效減少 dependencies

SYSDDIR="$GOPATH0/src/github.com/hacking-thursday/sysd"

function sync_sysd_dir(){
    if [ -d "$SYSDDIR" -o -L "$SYSDDIR" ] ;then
        rm -rf "$SYSDDIR"
    fi

    echo "Syncing files: $ROOT => $SYSDDIR "
    mkdir -p "$SYSDDIR"
    cp -r "$ROOT"/* "$SYSDDIR/"
}

if [ ! -d "$TMPDIR" ]; then mkdir -p "$TMPDIR" ; fi

sync_sysd_dir

pushd "$ROOT/sysd"
    go build -v -tags "$BuildTags"
popd
