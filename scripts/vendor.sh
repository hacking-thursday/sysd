#!/usr/bin/env bash

set -e

ROOT="$( cd $( dirname $0 ) && pwd -P )/.."
VENDOR="$ROOT/vendor"

# Downloads dependencies into vendor/ directory
cd $ROOT
test -d "$VENDOR" && rm -rf "$VENDOR"
mkdir -p "$VENDOR"
cd "$VENDOR"

clone() {
	vcs=$1
	rev=$2
	pkg=$3
	
	pkg_url=https://$pkg
	target_dir=src/$pkg
	
	echo -n "$pkg @ $rev: "
	
	if [ -d $target_dir ]; then
		echo -n 'rm old, '
		rm -fr $target_dir
	fi
	
	echo -n 'clone, '
	case $vcs in
		git)
			git clone --quiet --no-checkout $pkg_url $target_dir
			( cd $target_dir && git reset --quiet --hard $rev )
                        cur_rev=$( cd $target_dir; git log --oneline | head -1 | awk '{print $1}' )
                        echo -n " cur_rev=$cur_rev "
			;;
		hg)
			hg clone --quiet --updaterev $rev $pkg_url $target_dir
                        cur_rev=$( cd $target_dir; hg par | grep -e 'changeset' | awk '{print $2}' | cut -d: -f2 )
                        echo -n " cur_rev=$cur_rev "
			;;
	esac
	
	echo -n 'rm VCS, '
	( cd $target_dir && rm -rf .{git,hg} )
	
	echo done
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

clone git f92b795 github.com/Sirupsen/logrus
clone git e444e69 github.com/gorilla/mux
clone git 14f550f github.com/gorilla/context
clone git 6070b2c github.com/tsaikd/KDGoLib
clone git 3afe9db github.com/docker/docker

do_patch "$ROOT/misc/001.patch" "$VENDOR/src/github.com/docker/libcontainer"
do_patch "$ROOT/misc/002.patch" "$VENDOR/src/github.com/docker/docker"
do_patch "$ROOT/misc/003.patch" "$VENDOR/src/github.com/docker/docker"

rm -r "$VENDOR/src/github.com/docker/docker/docs"
rm -r "$VENDOR/src/github.com/docker/docker/vendor"
