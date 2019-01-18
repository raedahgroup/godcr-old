# godcr

## Overview
**godcr** is a [decred](https://www.decred.org/) wallet application for Desktop Operating Systems (Linux, macOS, Windows etc).
[dcrlibwallet](https://github.com/raedahgroup/dcrlibwallet/tree/dcrlibwallet-wip), a standalone decred wallet library, is used for all wallet access and control functionality.
**godcr** can also interface with [dcrwallet](https://github.com/decred/dcrwallet) over gRPC as an alternative to dcrlibwallet.

## Requirements
You can run **godcr** without installing any other software.

However, to use dcrwallet instead of dcrlibwallet for wallet operations, you'll need a running dcrwallet daemon.
Follow the steps below to download, setup and run dcrwallet:

* Download the **decred** release binaries for your operating system from [here](https://github.com/decred/decred-binaries/releases). Check under **Assets**.
* By default, **dcrwallet** uses **dcrd** to connect to the Decred network. The decred archive downloaded from the release page contains both binaries.
* After downloading and extracting **dcrd** and **dcrwallet**, [go here](https://docs.decred.org/wallets/cli/cli-installation/) to learn how to setup and run both binaries.

## Installation

### Option 1: Get the binary
**godcr** is not released yet. This doc will be updated with links to download the godcr binary when a release is ready. For now, build from source.

### Option 2: Build from source

#### Step 1. Install Go
* Minimum supported version is 1.11.4. Installation instructions can be found [here](https://golang.org/doc/install).
* **Ensure** `$GOPATH` environment variable is set and `$GOPATH/bin` is added to your PATH environment variable as part of the go installation process.

#### Step 2a. Install [QT](https://en.wikipedia.org/wiki/Qt_(software)) Binding for Go
```bash
go get -u -v github.com/therecipe/qt/cmd/...
```
If you get `module source tree too big` error message, try the following work around:
```bash
git clone https://github.com/therecipe/qt $GOPATH/src/github.com/therecipe/qt
cd $GOPATH/src/github.com/therecipe/qt/cmd
go install ./qtsetup
go install ./qtmoc
go install ./qtrcc
go install ./qtminimal
go install ./qtdeploy
```

#### Step 2b. Setup QT Binding for Go (Linux and Mac)
Run the following with `GO111MODULE=off` and outside $GOPATH
```bash
qtsetup test && qtsetup
```
It may be necessary to install additional dependencies on Linux. See [here](https://github.com/therecipe/qt/wiki/Installation-on-Linux).

#### Step 2b. Setup QT Binding for Go (Windows)
```bash
for /f %v in ('go env GOPATH') do %v\bin\qtsetup test && %v\bin\qtsetup
```

If building on Windows, there are additional steps to take to be able to build successfully. Those steps are described in the [setup instructions for Windows](https://github.com/therecipe/qt/wiki/Installation-on-Windows#if-you-want-to-install-the-binding), under **If you want to install the binding**

#### Step 2c. Detailed Installation Instructions
If you have issues with Step 2a or 2b above, you might need to consult the detailed setup instructions for [Windows](https://github.com/therecipe/qt/wiki/Installation-on-Windows), [Linux](https://github.com/therecipe/qt/wiki/Installation-on-Linux) or [MacOS](https://github.com/therecipe/qt/wiki/Installation-on-macOS). Focus only on the steps listed in the **Fast track version** section.

#### Step 3. Clone this repo
```bash
git clone https://github.com/raedahgroup/godcr $GOPATH/src/github.com/raedahgroup/godcr
```
**Note:** Cloning to a different directory other than `$GOPATH/src/github.com/raedahgroup` may cause build issues.

#### Step 4. Build the source code
```bash
cd $GOPATH/src/github.com/raedahgroup/godcr
export GO111MODULE=on
go mod download
go mod vendor
export GO111MODULE=off
go build
```
**Notes**
- Go modules must be enabled first to download all dependencies listed in `go.mod` to `vendor` folder within the project directory.
- Go modules must be disabled before running `go build` else the build will fail.
- In Windows, command prompt should always be restarted after changing environment variables for the changes to take effect.
- If you get checksum mismatch error while downloading dependencies**, ensure you're on go version 1.11.4 or higher and clean your go mod cache by running `go clean -modcache`

## Running godcr
### General usage
By default, **godcr** runs as a [cli app](https://en.wikipedia.org/wiki/Command-line_interface) where various wallet operations are performed by issuing commands on the terminal in the format:
```bash
godcr [options] <command> [args]
```
- Run `godcr -h` or `godcr help` to get general information of commands and options that can be issued on the cli.
- Use `godcr <command> -h` or   `godcr help <command>` to get detailed information about a command.

### As a GUI app
**godcr** can also be run as a full [GUI app](https://en.wikipedia.org/wiki/Graphical_user_interface) where wallet operations are performed by interacting with a graphical user interface.
The following GUI interface modes are supported:
1. Full GUI app on terminal.
Run `godcr --mode=terminal`
2. Web app served over http or https.
Run `godcr --mode=http`
3. Native desktop app with [nuklear](https://github.com/aarzilli/nucular) library.
Run `godcr --mode=nuklear`
4. Native desktop app with [qt](https://github.com/therecipe/qt) library.
Run`godcr --mode=qt`

### Configuration
The behaviour of the godcr program can be customized by editing the godcr configuration file.
The config file is where you set most options used by the godcr app, such as:
- the host and port to use for the http web server (if running godcr with `--mode=http`)
- the default interface mode to run (if you're tired of having to set `--mode=` everytime you run godcr)
- whether or not to use dcrwallet over gRPC for wallet functionality

Run `godcr -h` to see the location of the config file. Open the file with a text editor to see all customizable options.

## Contributing 

See the CONTRIBUTING.md file for details. Here's an overview:

1. Fork this repo to your github account
2. Before starting any work, ensure the master branch of your forked repo is even with this repo's master branch
2. Create a branch for your work (`git checkout -b my-work master`)
3. Write your codes
4. Commit and push to the newly created branch on your forked repo
5. Create a [pull request](https://github.com/raedahgroup/godcr/pulls) from your new branch to this repo's master branch
