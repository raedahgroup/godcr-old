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
* Set `$GOPATH` environment variable and add `$GOPATH/bin` to your PATH environment variable as part of the go installation process.

#### Step 2. Clone this repo
It is conventional to clone to $GOPATH, but not necessary.

**Linux**
```bash
git clone https://github.com/raedahgroup/godcr $GOPATH/src/github.com/raedahgroup/godcr
```

**Windows**
```
git clone https://github.com/raedahgroup/godcr %GOPATH%/src/github.com/raedahgroup/godcr
```

#### Step 3. Build the source code
* If you cloned to $GOPATH, set the `GO111MODULE=on` environment variable before building.
Run `export GO111MODULE=on` in terminal (for Mac/Linux) or `setx GO111MODULE on` in command prompt for Windows.
* `cd` to the cloned project directory.
* To build/install godcr-fyne:
`go build ./cmd/godcr-fyne` or `go install ./cmd/godcr-fyne`.
* To build/install the binaries for other interfaces:
`cd ./cmd/godcr-{interface} && go build` or `cd ./cmd/godcr-{interface} && go install`.
* Currently supported interfaces are `godcr-cli`, `godcr-fyne`, `godcr-nuklear`, `godcr-terminal` and `godcr-web`. 
* To run the http interface (`godcr-web`), you'd need to also build the frontend assets:
`cd ./web/static/app && yarn install && yarn build`.
You can get yarn from [here](https://yarnpkg.com/lang/en/docs/install/)

**Note: Building on Windows**
Exporting `GO111MODULE` directly in CLI does not work and
it is recommended to trigger via a `.bat` file.

* Create `enablegomod.bat` in the cloned project directory
* Paste:
  ```
  setx GO111MODULE on
  go mod download
  ```
* Execute `./enablegomod.bat`

## Running godcr
### As a CLI app
`godcr-cli` runs the program in a [command-line interface](https://en.wikipedia.org/wiki/Command-line_interface)
where various wallet operations are performed by issuing commands on the terminal in the format:
```bash
godcr-cli [options] <command> [args]
```
- Run `godcr-cli -h` or `godcr-cli help` to get general information of commands and options that can be issued on the cli.
- Use `godcr-cli <command> -h` or   `godcr-cli help <command>` to get detailed information about a command.

### As a GUI app
**godcr** can also be run as a full [GUI app](https://en.wikipedia.org/wiki/Graphical_user_interface)
where wallet operations are performed by interacting with a graphical user interface.
The following GUI interface modes are supported:
1. Full GUI app on terminal.
Run `godcr-terminal`.
2. Web app served over http or https.
Run `godcr-web`,
3. Native desktop app with [nuklear](https://github.com/aarzilli/nucular) library.
Run `godcr-nuklear`.
4. Native desktop app with [fyne](https://github.com/fyne-io/fyne) library.
Run `godcr-fyne`.

### Configuration
The behaviour of the godcr program can be customized by editing the godcr configuration file.
The config file is where you set most options used by the godcr app, such as:
- the host and port to use for the http web server (if running `godcr-web`)
- whether or not to use dcrwallet over gRPC for wallet functionality. 
To use dcrwallet, set the dcrwallet rpc address in config (e.g. `wallerrpcaddress=localhost:19111`).
If the rpc address is set in config, connection is made to dcrwallet. If not, dcrlibwallet is used.

Run `godcr-cli -h` to see the location of the config file.
Open the file with a text editor to see all customizable options.

### Features
[Go here](status.md) to view updated information about implemented features and known issues and workarounds.

## Contributing

See the CONTRIBUTING.md file for details. Here's an overview:

1. Fork this repo to your github account
2. Before starting any work, ensure the master branch of your forked repo is even with this repo's master branch
2. Create a branch for your work (`git checkout -b my-work master`)
3. Write your codes
4. Commit and push to the newly created branch on your forked repo
5. Create a [pull request](https://github.com/raedahgroup/godcr/pulls) from your new branch to this repo's master branch
