# Godcr Status Report

The current state of godcr is reported below.
Working features, known bugs and issues, work-in-progress features are all listed.

This document can also serve as a user's manual, showing how godcr works and how the various supported features can be accessed.

## Running Godcr
### General usage
By default, **godcr** runs as a [cli app](https://en.wikipedia.org/wiki/Command-line_interface) where various wallet operations are performed by issuing commands on the terminal in the format:
```bash
godcr [options] <command> [args]
```
- Run `godcr -h` or `godcr help` to get general information of commands and options that can be issued on the cli.
- Use `godcr <command> -h` or   `godcr help <command>` to get detailed information about a command.

### Godcr GUI
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

## Godcr Features
The following table lists all features Godcr currently supports, and the level to which each feature is supported on each interface.

| Feature | cli | http | nuklear | qt | tview
|---|---|---|---|---|---|
| create wallet | :white_check_mark: |
| blockchain sync | :white_check_mark: |
| balance | :white_check_mark: |
| receive | :white_check_mark: |
| send | :white_check_mark: |
| send custom | :white_check_mark: |
| history | :white_check_mark: |

**For cli**

| Operation | Status | Notes
|---|---|---|
| `dcrcli -v` | :white_check_mark: |
| `dcrcli [command] -h` | :white_check_mark: |
| `dcrcli --createwallet` | :white_check_mark: |
| `dcrcli [command] --sync` | :white_check_mark: |
| `dcrcli balance [-d/--detailed]` | :white_check_mark: |
| `dcrcli createaccount <accountName>` _(experimental feature)_ | :white_check_mark: |
| `dcrcli receive [account-name]` | :white_check_mark: |
| `dcrcli send` | :white_check_mark: | **requires** `--sync` flag when using `mobilewallet` |
| `dcrcli send-custom` | :white_check_mark: | **requires** `--sync` flag when using `mobilewallet` |
| `dcrcli history` | :white_check_mark: |

**For web**
