#!/usr/bin/env bash

set -e

ROOT="$( cd $( dirname $0 ) && pwd -P )"

export GOPATH0="$ROOT/.gopath" 
if [ -n "$GOPATH" ];then
    # 若有預設的 GOPATH，就優先使用預設的 GOPATH 為主
    export GOPATH="$GOPATH:$GOPATH0" 
else
    export GOPATH="$GOPATH0" 
fi
export GOPATH9="$( echo $GOPATH | cut -d: -f1)" 
echo "GOPATH: $GOPATH"
echo "GOPATH0: $GOPATH0"
echo "GOPATH9: $GOPATH9"

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

function do_patch(){
    set +e

    local PATCH_FILE="$1"
    local TARGET_DIR="$2"

    PATCH_OPTION="--force --batch -p1"
    if [ -d "$TARGET_DIR" ]; then
        patch --dry-run $PATCH_OPTION -d "$TARGET_DIR" < "$PATCH_FILE" > /dev/null
        if [ $? -eq 0 ]; then
            patch $PATCH_OPTION -d "$TARGET_DIR" < "$PATCH_FILE"
        fi
    fi

    set -e
}

if [ ! -d "$TMPDIR" ]; then mkdir -p "$TMPDIR" ; fi
if [ ! -d "$SYSDDIR" ]; then mkdir -p "$SYSDDIR" ; fi

sync_sysd_dir
pushd "$SYSDDIR"
    do_patch "$ROOT/misc/001.patch" "${GOPATH9}/src/github.com/docker/libcontainer"
    do_patch "$ROOT/misc/002.patch" "${GOPATH9}/src/github.com/docker/docker"
    do_patch "$ROOT/misc/003.patch" "${GOPATH9}/src/github.com/docker/docker"
    go get -v -t -tags "$BuildTags" ./sysd
    if [ $? -eq 0 ];then
        PASS_DEPS="ok"
    fi
popd

if [ "$PASS_DEPS" = "ok" ];then
    # build first
    pushd "$ROOT/sysd"
        go build -v -tags "$BuildTags"
        # if [ $? -eq 0 ];then
        #     # test later
        #     pushd "$SYSDDIR"
        #         go test -v -tags "$BuildTags" ./...
        #     popd
        # fi
    popd
fi
