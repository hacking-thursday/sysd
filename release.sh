#!/usr/bin/env bash

get_latest_ver(){
    git tag --sort=-refname | head -1 | cut -c2-
}

make_tarball(){
    local ver="$1"
    local out="$2"

    git archive --format=tar.gz --output=$out --prefix=sysd-$ver/ "v${ver}"
}

make_deb_pkg(){
    local tarball="$1"

    echo $tarball
    local srcdir=$( tar -tf $tarball | head -1)

    tar -xzvf $tarball
    if [ -d  "$srcdir" ];then
        cp -r debian/ "$srcdir"
        tar -zcv --exclude="./bin" --exclude="./pkg" --exclude=".git" \
                  --exclude="docker/vendor" \
                  --exclude="docker/docs" \
                  --exclude="docker/hack" \
                  --exclude="docker/contrib" \
                  --exclude="*/testdata/*" \
                  --exclude="hacking-thursday/sysd" \
                  -C .gopath -f "/tmp/deps.tar.gz" .
        install -d "$srcdir/.gopath" 
        tar -zxvf "/tmp/deps.tar.gz" -C "$srcdir/.gopath"
        pushd "$srcdir"
            debuild -S
        popd
    fi
}

version=$( get_latest_ver )

make_tarball "$version" "sysd_$version.orig.tar.gz"

make_deb_pkg "sysd_$version.orig.tar.gz"
