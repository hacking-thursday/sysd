#!/usr/bin/env bash

set -e

ROOT=$( readlink -f $( dirname $0 ) )
export GOPATH="$ROOT/.gopath" 
export TMPDIR="$ROOT/.tmp" 


H4DIR="$GOPATH/src/github.com/hacking-thursday"
SYSDDIR="$H4DIR/sysd"

function replace_sysd_dir(){
    pushd $H4DIR > /dev/null
        if [ -d "./sysd" ]; then
            rm -rvf "./sysd"
            ln -s ../../../../ ./sysd
        fi
    popd > /dev/null
}

if [ ! -d $TMPDIR ]; then mkdir -p $TMPDIR ; fi
if [ ! -d $SYSDDIR ]; then mkdir -p $SYSDDIR ; fi

if [ -d $SYSDDIR -a ! -L $SYSDDIR ] ;then
    cp -vr $ROOT/* $SYSDDIR/
    pushd $SYSDDIR
        go get -v -t ./...
        go test ./...
        if [ $? -eq 0 ];then
            replace_sysd_dir
        fi
    popd
fi

if [ -L $SYSDDIR ];then
    pushd $SYSDDIR
        cd ./sysd; go build -v
    popd
fi
