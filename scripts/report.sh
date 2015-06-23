#!/usr/bin/env bash

DATE="2014-11-06"

function convert_author_map(){
    sed -e 's/matlinuxer2 <matlinuxer2@gmail.com>/Chun-Yu Lee (Mat) <matlinuxer2@gmail.com>/g'
}

function list_authors(){
cat <<EOF
matlinuxer2
Carl
Bruce
yan
tsaikd
EOF
}

function sort_cnt(){
    sort | uniq -c | sort -k1 -n
}

mkdir reports
git log --stat --since="$DATE"  --author=".*" | grep Author | convert_author_map | sort_cnt > reports/$DATE-ALL.log
for uu in `list_authors`; do 
  git log --stat --since="$DATE"  --author="$uu" > reports/$DATE-$uu.log
done
