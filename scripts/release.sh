#!/usr/bin/env bash

ROOT="$( cd $( dirname $0 ) && pwd -P )/.."

function check_version_format(){
    local ver="$1"

    echo "$ver" | grep -e '^v[0-9]\+\.[0-9]\+\.[0-9]\+$' 2>&1 >/dev/null
    
    return $?
}

function get_previous_version_from_changelog(){
    local changelog_path="$1"

    rev=$( grep    -e   '^## \([0-9]\+\.[0-9]\+\.[0-9]\+\) (.*).*' "$changelog_path" \
      | sed -e 's/^## \([0-9]\+\.[0-9]\+\.[0-9]\+\) (.*).*/\1/g' | tail -n +1 | head -1 )

    echo "v$rev"
}

VER="$1"
VER2=${VER#v}

check_version_format $VER || ( echo "版本編號格式不佳, 範例: v0.6.4" ; exit 1 )

pushd $ROOT > /dev/null
    # 1. 修改 VERSION
    echo $VER2 > ./VERSION

    # 2. 檢查 Changelog
    cat ./CHANGELOG.md | grep -e "^##\s*$VER2\s*([0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9])\s*$"
    if [ $? -ne 0 ];then
        echo "尚未更新 CHANGELOG.md" 

        VER_prev=$( get_previous_version_from_changelog "./CHANGELOG.md" )
        echo "CHANGELOG.md 上一個版號是: $VER_prev"

        echo "初步自動產生 $VER_prev ~ $VER 的 commit log 進 CHANGELOG.md"
        tmpf=$( mktemp )
        date_today=$( date +"%Y-%m-%d" )
        echo "## $VER2 ($date_today)" >> $tmpf
        echo "" >> $tmpf
        git diff --name-status ${VER_prev}..HEAD >> $tmpf
        echo "" >> $tmpf
        echo "Features:" >> $tmpf
        echo "" >> $tmpf
        git log --no-merges --pretty=format:'  - %s ( by %an )' ${VER_prev}..HEAD >> $tmpf
        echo "" >> $tmpf
        echo "" >> $tmpf
        echo "Bugfixes:" >> $tmpf
        echo "" >> $tmpf
        cat ./CHANGELOG.md >> $tmpf
        cat $tmpf > ./CHANGELOG.md
        rm -f $tmpf
    fi

    # 3. 更新 Manifest
    echo "更新 Manifest"
    make Manifest

    # 4. 製作 tarball
    echo "製作 tarball"
    make dist

    # 5. 簡單測試 tarball
    echo "簡單測試 tarball"
    tgz="sysd-${VER2}.tar.gz"
    tgzdir="sysd-${VER2}/"
    if [ -e  "$tgz" ]; then
        tar -zxvf $tgz
        pushd $tgzdir
            ./configure && make 
        popd
    fi
    rm -r $tgzdir
popd > /dev/null
