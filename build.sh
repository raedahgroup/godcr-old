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

function deploy() {
  buildFrontEnd
  echo 'go get -u github.com/gobuffalo/packr/v2/packr2'
  go get -u github.com/gobuffalo/packr/v2/packr2
  echo 'go mod download'
  go mod download
  echo 'packr build'
  packr build
}
ACTION=$1
if [[ "$ACTION" = "web" ]]; then
  buildBackEnd
  buildFrontEnd
else
  deploy
fi