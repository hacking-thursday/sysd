#!/usr/bin/env bash

ROOT="$( cd $( dirname $0 ) && pwd -P )/../.."

get_gpg_signname(){
        local result=$( gpg --list-secret-keys| grep uid | head --lines=1 | cut -c 22- )
        echo $result
}

function local_dput(){
	local TARGET="$1"

	local DPUT_CONF="/tmp/._dput.cf"

	cat > $DPUT_CONF <<EOD
[ppaa]
fqdn                    = ppa.launchpad.net
method                  = ftp
incoming                = ~sysd/sysd
login                   = anonymous
EOD

	dput -c $DPUT_CONF ppaa $TARGET

	rm $DPUT_CONF
}

PKG_NAME="sysd"
PKG_VER="0.6.0"
CNT="1"
TARBALL_URL="https://github.com/hacking-thursday/sysd/releases/download/v${PKG_VER}/sysd-${PKG_VER}.tar.gz"
TARBALL_NAME="${PKG_NAME}-${PKG_VER}.tar.gz"
TARBALL_NAME2="${PKG_NAME}_${PKG_VER}.orig.tar.gz"
TARBALL_DIR="${TARBALL_NAME%.tar.gz}"
TEMP_DIR="__debian__$(date '+%Y-%m-%d_%H%M%S')"

if [ ! -f "$ROOT/pkg/tgz/$TARBALL_NAME" ];then
    wget -O "$ROOT/pkg/tgz/$TARBALL_NAME" "$TARBALL_URL"
fi


mkdir $TEMP_DIR
pushd $TEMP_DIR
	cp "$ROOT/pkg/tgz/$TARBALL_NAME" "$TARBALL_NAME2"
        tar -zxvf "$TARBALL_NAME2"
	cp -av ../debian "$TARBALL_DIR"

        for codename in "lucid" "precise" "trusty" "utopic" ; do
            cp "$TARBALL_DIR/debian/changelog" /tmp/changelog.bak
            pushd "$TARBALL_DIR"
                sed -i debian/control -e 's/golang-go,.*/golang-go/g'
                DEBEMAIL="$(get_gpg_signname)" debchange --newversion "${PKG_VER}-0" --force-distribution --distribution $codename "for $codename"
                debuild -S | tee /tmp/debuild.log

                CHANGES="$(sed -n 's/^.*signfile \(.*\.changes\).*$/\1/p' /tmp/debuild.log)"
                rm /tmp/debuild.log
                local_dput ../$CHANGES
            popd

            pushd "$TARBALL_DIR"
                for (( i=1; i<=$CNT; i++ )); do
                    DEBEMAIL="$(get_gpg_signname)" debchange --local=$codename --force-distribution --distribution $codename "for $codename"
                done
                debuild -S | tee /tmp/debuild.log

                CHANGES="$(sed -n 's/^.*signfile \(.*\.changes\).*$/\1/p' /tmp/debuild.log)"
                rm /tmp/debuild.log
                local_dput ../$CHANGES
            popd
            cp /tmp/changelog.bak "$TARBALL_DIR/debian/changelog"
        done
popd
