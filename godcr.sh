#!/usr/bin/env bash
function buildFe {
  echo "cd web/static/app"
  cd web/static/app
  echo 'yarn build'
  yarn build
  echo 'cd ../../../'
  cd ../../../
}
function buildBe {
  echo 'go build'
  go build
}
function deploy() {
  buildFe
  echo 'go mod download'
  go mod download
  echo 'packr build'
  packr build
}
ACTION=$1
if [[ "$ACTION" = "build" ]]; then
  buildBe
elif [[ "$ACTION" = "build-web" ]]; then
  buildFe
  buildBe
elif [[ "$ACTION" = "deploy" ]]; then
  deploy
else
  echo "./godcr $1 $2 $3 $4"
  ./godcr $1 $2 $3 $4
fi