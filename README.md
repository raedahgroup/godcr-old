# godcr
====

## Dcrcli Overview 

godcr is a command-line utility that interfaces with [Dcrwallet](https://github.com/decred/dcrwallet) rpc's methods.

## Requirements 
* [Go](http://golang.org) 1.11 
* Git
* Running `dcrwallet` instance 

## Installation 

### Build from source

The following guide assumes a Unix-like shell (e.g bash).

* [Install Go](http://golang.org/doc/install).
It is recommended to add `$GOPATH/bin` to your `PATH`.

* [Install Git](https://git-scm.com).

* Download dcrd and dcrwallet from [here](https://github.com/decred/decred-binaries/releases).
Click on assets and download the package for your operating system.

* After downloading the file, unzip or unpack the file then [follow this link](https://docs.decred.org/wallets/cli/cli-installation/)
 to install, setup and run drcd and dcrwallet on your machine.

* Clone this repo to somewhere on your computer. Set the `GO111MODULE=on` environment variable if you are building from within `GOPATH`.

##### Example of obtaining and building from source in Linux 
```bash 
$ git clone https://github.com/raedahgroup/godcr
$ cd godcr
$ go install or GO111MODULE=on go install (if you are building from within `GOPATH`)
```

## Running godcr 

### Create configuration file 

Begin with the sample configuration file:

```bash 
cp sample-godcr.conf godcr.conf 
``` 

Then edit godcr.conf and input your RPC settings. After you are finished, move godcr.conf to the `appdata` folder (default is `~/.godcr` on Linux, `%localappdata%\Godcr` on Windows). See the output of `godcr -h` for a list of all options.

### Using godcr

Run `godcr <command> <options>`. See the output of `godcr -l` for a list of all commands.

## Contributing 

See the CONTRIBUTING.md file for details. Here's an overview of it: 

1. Fork the repo
1. Create a branch for your work (`git branch -b branch`).
3. Write your codes 
4. Commit and push to your repo
5. Create a [pull request](https://github.com/raedahgroup/godcr)

## License
