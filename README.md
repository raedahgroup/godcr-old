# godcr

## Overview
**godcr** is a decred wallet application for Linux, macOS and Windows that provides wallet access and control functionality using [dcrlibwallet](https://github.com/raedahgroup/godcr/pull/88). It can also interface with [dcrwallet](https://github.com/decred/dcrwallet) via RPC as an alternative to dcrlibwallet. The godcr app can be run in any of the following interface modes:
- Web (web app running on an http server)
- Cli (command-line utility)
- Desktop (native desktop application, currently work in progress)

## Requirements
You can run **godcr** without installing any other software. By default **godcr** uses dcrlibwallet.
To run **godcr** using **dcrwallet** instead of [dcrlibwallet](https://github.com/raedahgroup/godcr/pull/88), the following is required.
* Download the **decred** release binaries for your operating system from [here](https://github.com/decred/decred-binaries/releases). Check under **Assets**.
* **dcrwallet** requires **dcrd** to work. The decred archive downloaded from the release page contains both binaries.
* After downloading and extracting **dcrd** and **dcrwallet**, [go here](https://docs.decred.org/wallets/cli/cli-installation/) to learn how to setup and run both binaries.

## Installation

### Option 1: Get the binary
**godcr** is not released yet. This doc will be updated with links to download the godcr binary when a release is ready. For now, build from source.

### Option 2: Build from source
* Install Go. Installation instructions can be found [here](https://golang.org/doc/install). It is recommended to add $GOPATH/bin to your PATH as part of the installation process.
* Clone this repository. It is conventional to clone to $GOPATH, but not necessary.
```bash
$ git clone https://github.com/raedahgroup/godcr $GOPATH/src/github.com/raedahgroup/godcr
```
* If you clone to $GOPATH, set the `GO111MODULE=on` environment variable when building. On Unix systems, you can add the following line to `~/.bash_profile` to persist the variable
```bash
export GO111MODULE=on
```
* `cd` to the cloned project directory and build or install godcr. Building will place the godcr binary in your working directory while install will place the binary in $GOPATH/bin
```bash
$ go build
$ go install
```
* If you cloned the source code to $GOPATH but have not set the GO111MODULE=on environment variable, build/install like this
```bash
$ GO111MODULE=on go build
$ GO111MODULE=on go install
```

## Running godcr
* Create the godcr configuration file by copying the sample configuration file (`sample-godcr.conf`) to the default config location and editing as necessary. The default config location can be gotten by running `godcr -h`.
The settings in the configuration file can also be provided on the command line as options to the program. Flags passed on the command line override those set in the configuration file. Run `godcr -h` to see the available options.
* Once you've configured godcr, you can perform various wallet-related operations by running
```bash
$ godcr [options] <command> [args]
```
* To view available commands and options at any time:
```bash
$ godcr -h
```
* To get detailed help information for a command:
```bash
$ godcr <command> -h
```

## Contributing 

See the CONTRIBUTING.md file for details. Here's an overview:

1. Fork this repo to your github account
2. Before starting any work, ensure the master branch of your forked repo is even with this repo's master branch
2. Create a branch for your work (`git checkout -b my-work master`)
3. Write your codes
4. Commit and push to the newly created branch on your forked repo
5. Create a [pull request](https://github.com/raedahgroup/godcr/pulls) from your new branch to this repo's master branch
