#!/usr/bin/env bash
function deployWeb() {
   echo "building frontend assets with yarn"
   (cd web/static/app && yarn install && yarn build)
   echo "building godcr-web binary with packr2"
   (cd cmd/godcr-web && packr2 build && cd ../../)
   mv ./cmd/godcr-web/godcr-web ./godcr-web
   echo "binary saved to ./godcr-web"
}

function buildFyne() {
    echo "packing assets with packr2"
    (cd fyne && packr2)
    echo "building with go build"
    cd cmd/godcr-fyne && go build && cd ../../
   mv ./cmd/godcr-web/godcr-fyne ./godcr-fyne
   echo "binary saved to ./godcr-fyne"
}

function buildNuklear() {
    echo "building with go build"
    cd ./cmd/godcr-nuklear && go build && cd ../../
   mv ./cmd/godcr-web/godcr-nuklear ./godcr-nuklear
   echo "binary saved to ./godcr-nuklear"
}

function buildTerminal() {
    echo "building with go build"
    cd ./cmd/godcr-terminal && go build && cd ../../
   mv ./cmd/godcr-web/godcr-terminal ./godcr-terminal
   echo "binary saved to ./godcr-terminal"
}

function buildCli() {
    echo "building with go build"
    cd ../ && go build ./cmd/godcr-cli && cd ../../
   mv ./cmd/godcr-web/godcr-cli ./godcr-cli
   echo "binary saved to ./godcr-cli"
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
else
    echo "Usage: ./build.sh {interface} e.g. ./build.sh web"
fi
