#!/usr/bin/env bash
function buildFrontEnd {
  echo "cd web/static/app"
  cd web/static/app
  echo 'yarn install'
  yarn install
  echo 'yarn build'
  yarn build
  echo 'cd ../../../'
  cd ../../../
}

function buildBackEnd {
  echo 'go build'
  go build
}

function deployWeb() {
  buildFrontEnd
  cd cmd/godcr-web
  echo 'go mod download'
  go mod download
  echo 'packr2 build'
  packr2 build
  cd ../../
}
ACTION=$1
if [[ "$ACTION" = "web" ]]; then
  buildBackEnd
  buildFrontEnd
else
  deployWeb
fi