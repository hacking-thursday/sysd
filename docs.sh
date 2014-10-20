#!/usr/bin/env bash

SYSD_WIKI_URL="http://github.com/hacking-thursday/sysd/wiki"
SYSD_WIKI_DIR="./sysd.wiki"

function list_files(){
    find mods/ -type f | xargs grep Register | cut -d: -f1 | sort | uniq | grep -v '_test'
}

function parse_file(){
    local fpath="$1"

    params=$( cat "$fpath" |grep -e 'mods\.Register' | sed -e 's/.*(\(.*\)).*/\1/g' )
    name=$( cat "$fpath" | head | grep -e 'package\s\+' | awk '{print $2}' )

    if [ -n "$params" ];then
        echo -n "<tr>"
        echo -n "<td>" $( echo $params | cut -d, -f2 | sed -e 's/"//g' -e "s/'//g" )  "</td>"
        echo -n "<td>" $( echo $params | cut -d, -f1 | sed -e 's/"//g' -e "s/'//g" )  "</td>"
        echo -n "<td>" "<a href=\"$SYSD_WIKI_URL/$name\">$name</a>" "</td>"
        echo -n "<td>" "$fpath" "</td>"
        echo -n "<td>" "?" "</td>"
        echo -n "<td>" "?" "</td>"
        echo -n "<td>" "?" "</td>"
        echo -n "</tr>"
        echo -n -e "\n"
    fi
}

function gen_api_listing(){
    echo ""
    echo "<!-- API Listing beg -->"
    echo ""
    echo "<table>"
    echo " <tr><th> API   </th><th> Method </th><th> Module </th><th>src path </th><th> Linux </th><th> Windows </th><th> OSX  </th> </tr>"
    for ff in $( list_files ); do 
        parse_file $ff
    done
    echo "<table>"
    echo ""
    echo "<!-- API Listing end -->"
    echo ""
}


if [ ! -d "$SYSD_WIKI_DIR" ];then
    echo "Please checkout sysd's wiki first:"
    echo ""
    echo "      git clone git@github.com:hacking-thursday/sysd.wiki.git \"$SYSD_WIKI_DIR\" "
    echo ""
    exit 0
fi


if [ -d "$SYSD_WIKI_DIR" ];then
    cat "$SYSD_WIKI_DIR/API.md" | grep -e ' API Listing ' 
    if [ $? -eq 0 ];then
        cat "$SYSD_WIKI_DIR/API.md" | grep -e '<!-- API Listing beg -->' -B 9999 | head -n -2 > /tmp/head.txt
        cat "$SYSD_WIKI_DIR/API.md" | grep -e '<!-- API Listing end -->' -A 9999 | tail -n +3 > /tmp/tail.txt
    else 
        cat "$SYSD_WIKI_DIR/API.md"  > /tmp/head.txt
        echo "" > /tmp/tail.txt 
    fi
    gen_api_listing > /tmp/middle.txt

    cat /tmp/head.txt /tmp/middle.txt /tmp/tail.txt > "$SYSD_WIKI_DIR/API.md"
fi
