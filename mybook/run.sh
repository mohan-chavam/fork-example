#!/usr/bin/env bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
end='\033[0m'

BINARY_NAME=mybook.bin
install_app() {
    cd mybook-static &&
        npm run build &&
        cp ../../toolbox/keylogger/static/favicon.ico dist/favicon.ico &&
        cd .. &&
        statik -f -src=mybook-static/dist/ -dest app/common/ &&
        go build -o bin/${BINARY_NAME}

    if [ ! -d bin/data ]; then
        ln -s data bin/data
    fi

    bin/${BINARY_NAME} -s -p 9090
}

install_server() {
    go build -o bin/${BINARY_NAME}

    if [ ! -d bin/data ]; then
        ln -s data bin/data
    fi

    bin/${BINARY_NAME} -s -p 9090
}

run(){
  bin/${BINARY_NAME} -s -p 9090
}

help() {
    printf "Run：$red sh $0 $green<verb> $yellow<args>$end\n"
    format="  $green%-6s $yellow%-8s$end%-20s\n"
    printf "$format" "-h" "" "帮助"
}

case $1 in
-h)
    help
    ;;
-r)
    run
    ;;
-s)
    install_server
    ;;
*)
    install_app
    ;;
esac
