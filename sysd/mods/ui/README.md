## Build static frontend

NOTE: node, npm is required for compilation

    pushd $(git rev-parse --show-toplevel)/mods/ui/
    npm install gulp
    export PATH="$PATH:`pwd`/node_modules/.bin/"
    gulp build
    popd
