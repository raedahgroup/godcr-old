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

## Godcr Features (Summary)
| Feature | Description | Status | Next Steps |
|---|---|---|---|
| create wallet | If no wallet exists, user is asked to create one using this feature. | cli :white_check_mark: <br> http :white_check_mark: <br> nuklear :white_check_mark: <br> qt :white_check_mark: <br> tview :white_check_mark: | Allow creating multiple wallets, even if wallet already exists |
| sync blockchain (spv) | Blockchain sync is performed everytime godcr is launched in gui mode.<br>In cli mode, the `--sync` flag is used to trigger a blockchain sync operation. | cli :white_check_mark: <br> http :white_check_mark: <br> nuklear :white_check_mark: <br> qt :white_check_mark: <br> tview :white_check_mark: | Allow creating multiple wallets, even if wallet already exists
| sync blockchain (rpc) | Similar to above feature, syncs blockchain by connecting to a running instance of dcrd over rpc | :x: | Support for this feature should be added to all interfaces |
| balance | Show balance for all accounts in wallet | :white_check_mark: |
| receive | Generate address to receive funds | :white_check_mark: |
| send funds (simple) | Send funds to 1 or more decred addresses | :white_check_mark: |
| send funds (custom) | Similar to above, with ability to customize inputs and change outputs | :white_check_mark: |
| history | View wallet transaction history | :white_check_mark: |
| tx detail | Show detailed information for any wallet transaction | :white_check_mark: |
| stake info | View status of purchased tickets and stake info | :white_check_mark: |
| purchase ticket(s) | Purchase 1 or more tickets | :white_check_mark: |

## Cli - Known Issues and Additional Information
#### Sync blockchain
- Unlike in other interfaces, the cli interface does not automatically sync the blockchain before performing wallet operations.
This may lead to inaccuracy of displayed information or complete inability to perform certain wallet operations such as `send` and `purchasetickets`
- To circumvent the issue(s) identified above, the `--sync` flag should be used when issuing godcr commands on cli. Alternatively, you can set `sync=1` or `sync=true` in `godcr.conf` to always perform a blockchain sync before performing any wallet operation.
- Also, `godcr --sync` can be run alone, without any command to perform a blockchin sync at any time.
- The above concerns do not apply if godcr is _properly_ configured to perform wallet operations using dcrwallet rpc.
