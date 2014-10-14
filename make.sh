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
echo "GOPATH: $GOPATH"

export TMPDIR="$ROOT/.tmp" 
export CGO_ENABLED="0" # 有效減少 dependencies
export BuildTags='exclude_graphdriver_devicemapper exclude_graphdriver_aufs exclude_graphdriver_btrfs' # 有效減少 dependencies

H4DIR="$GOPATH0/src/github.com/hacking-thursday"
SYSDDIR="$H4DIR/sysd"

function replace_sysd_dir(){
    pushd "$H4DIR" > /dev/null
        if [ -d "./sysd" ]; then
            rm -rf "./sysd"
            ln -s ../../../../ ./sysd
        fi
    popd > /dev/null
}

function do_patch(){
    set +e

    local PATCH_FILE="$1"
    local TARGET_DIR="$2"

    PATCH_OPTION="--force --batch -p1"
    patch --dry-run $PATCH_OPTION -d "$TARGET_DIR" < "$PATCH_FILE" > /dev/null
    if [ $? -eq 0 ]; then
        patch $PATCH_OPTION -d "$TARGET_DIR" < "$PATCH_FILE"
    fi

    set -e
}

if [ ! -d "$TMPDIR" ]; then mkdir -p "$TMPDIR" ; fi
if [ ! -d "$SYSDDIR" ]; then mkdir -p "$SYSDDIR" ; fi

if [ -L "$SYSDDIR" ];then 
    echo -n -e "[Symbolic link]\t"; 
    MESSAGE=" => $( cd $SYSDDIR && pwd -P )";
else
    echo -n -e "[Mirror copy]\t"; 
fi
echo "$SYSDDIR $MESSAGE"

if [ -d "$SYSDDIR" -a ! -L "$SYSDDIR" ] ;then
    cp -r "$ROOT"/* "$SYSDDIR/"
    pushd "$SYSDDIR"
        do_patch "$ROOT/misc/001.patch" "${GOPATH0}/src/github.com/docker/libcontainer"
        do_patch "$ROOT/misc/002.patch" "${GOPATH0}/src/github.com/docker/docker"
        go get -v -t -tags "$BuildTags" ./sysd
        if [ $? -eq 0 ];then
            replace_sysd_dir
        fi
    popd
fi

if [ -L "$SYSDDIR" ];then
    # build first
    pushd "$SYSDDIR/sysd"
        go build -v -tags "$BuildTags"
        # if [ $? -eq 0 ];then
        #     # test later
        #     pushd "$SYSDDIR"
        #         go test -v -tags "$BuildTags" ./...
        #     popd
        # fi
    popd
fi
