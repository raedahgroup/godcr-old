#!/usr/bin/env bash
function deployWeb() {
   echo "building frontend assets with yarn"
   (cd ./web/static/app && yarn install && yarn build)
   echo "building godcr-web with packr2"
   (cd ./cmd/godcr-web && packr2 build)
   mv ./cmd/godcr-web/godcr-web ./godcr-web
   echo "binary saved to ./godcr-web"
}

function buildFyne() {
    echo "packing assets with packr2"
    (cd fyne && packr2)
    echo "building godcr-fyne with go build"
    (cd ./cmd/godcr-fyne && go build)
    mv ./cmd/godcr-fyne/godcr-fyne ./godcr-fyne
    echo "binary saved to ./godcr-fyne"
}

function buildNuklear() {
    echo "building with go build"
    (cd ./cmd/godcr-nuklear && go build)
    mv ./cmd/godcr-nuklear/godcr-nuklear ./godcr-nuklear
    echo "binary saved to ./godcr-nuklear"
}

function buildTerminal() {
    echo "building with go build"
    (cd ./cmd/godcr-terminal && go build)
    mv ./cmd/godcr-terminal/godcr-terminal ./godcr-terminal
    echo "binary saved to ./godcr-terminal"
}

function buildCli() {
    echo "building with go build"
    (cd ./cmd/godcr-cli && go build)
    mv ./cmd/godcr-cli/godcr-cli ./godcr-cli
    echo "binary saved to ./godcr-cli"
}

function buildGio() {
    echo "building with go build"
    (cd ./cmd/godcr-gio && go build)
    mv ./cmd/godcr-gio/godcr-gio ./godcr-gio
    echo "binary saved to ./godcr-gio"
}

interface=$1
if [[ "$interface" = "web" ]]; then
    deployWeb
elif [[ "$interface" = "fyne" ]]; then
    buildFyne
elif [[ "$interface" = "nuklear" ]]; then
    buildNuklear
elif [[ "$interface" = "terminal" ]]; then
    buildTerminal
elif [[ "$interface" = "cli" ]]; then
    buildCli
elif [[ "$interface" == "gio" ]]; then 
    buildGio
else
    echo "Usage: ./build.sh {interface} e.g. ./build.sh web"
fi
