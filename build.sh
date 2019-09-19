#!/usr/bin/env bash
function buildWebFrontEnd {
  echo "cd web/static/app"
  cd web/static/app
  echo 'yarn install'
  yarn install
  echo 'yarn build'
  yarn build
  echo 'cd ../../../'
  cd ../../../
}

function deployWeb() {
  buildWebFrontEnd
  cd cmd/godcr-web
  echo 'go mod download'
  go mod download
  echo 'packr2 build'
  packr2 build
  cd ../../
}

function buildFyne() {
    echo "packing assets with packr2"
    (cd fyne && packr2)
    echo "building with go build"
    cd ../ && go build ./cmd/godcr-fyne
}

function buildNuklear() {
    echo "building with go build"
    cd ../ && go build ./cmd/godcr-nuklear
}

function buildTerminal() {
    echo "building with go build"
    cd ../ && go build ./cmd/godcr-terminal
}

function buildCli() {
    echo "building with go build"
    cd ../ && go build ./cmd/godcr-cli
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
