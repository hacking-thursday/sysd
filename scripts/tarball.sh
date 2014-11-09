#!/usr/bin/env bash

set -e

ROOT="$( cd $( dirname $0 ) && pwd -P )/.."
VENDOR="$ROOT/vendor"

VERSION=$( cat "$ROOT/VERSION" )
PKGNAME="sysd"

PKGDIR="$ROOT/$PKGNAME-$VERSION"
TGZDIR="$ROOT/pkg/tgz"

if [ -d "$PKGDIR" ];then
    rm -rf "$PKGDIR"
fi

mkdir -p "$PKGDIR"
mkdir -p "$TGZDIR"
pushd $ROOT > /dev/null
    for ff in `cat Manifest`;do
        cp -av --parents "$ff" "$PKGDIR"
    done
    tar -czf "$TGZDIR/$PKGNAME-$VERSION.tar.gz" "$PKGDIR"
popd > /dev/null

if [ -d "$PKGDIR" ];then
    rm -rf "$PKGDIR"
fi
