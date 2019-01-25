# Godcr Status
The current state of godcr is reported below.
Working features, known bugs and issues, work-in-progress features are all listed.

## Running Godcr
By default, **godcr** runs as a [cli app](https://en.wikipedia.org/wiki/Command-line_interface) where various wallet operations are performed by issuing commands on the terminal in the format:
```bash
godcr [options] <command> [args]
```
- Run `godcr -h` or `godcr help` to get general information of commands and options that can be issued on the cli.
- Use `godcr <command> -h` or   `godcr help <command>` to get detailed information about a command.

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
The list of features available on godcr is maintained below.
Some features are either only partially implemented (:ballot_box_with_check:) or not implemented at all (:x:) on some interfaces

| Feature | :white_check_mark: | :ballot_box_with_check: | :x: | Next Steps |
|---|---|---|---|---|
| create wallet | cli, http | terminal | nuklear, qt | Allow creating multiple wallets, even if wallet already exists<br><br>Ask for network type when creating wallet<br><br>Seed display confirmation should follow same pattern as dcrandroid |
| sync blockchain (spv) | cli, http | terminal | nuklear, qt |
| sync blockchain (rpc) | | | all | Add this feature to all interfaces |
| balance | cli, http | nuklear, qt | terminal |
| receive | cli, http, nuklear | | qt, terminal |
| send funds (simple) | cli, http | | nuklear, qt, terminal |
| send funds (custom) | cli | http | nuklear, qt, terminal |
| history | cli, http | | nuklear, qt, terminal |
| tx detail | cli, http | | nuklear, qt, terminal |
| stake info | cli, http | | nuklear, qt, terminal |
| purchase ticket(s) | cli, http | | nuklear, qt, terminal |

## Cli - Known Issues and Additional Information
#### Sync blockchain
- Unlike in other interfaces, the cli interface does not automatically sync the blockchain before performing wallet operations.
This may lead to inaccuracy of displayed information or complete inability to perform certain wallet operations such as `send` and `purchasetickets`
- To circumvent the issue(s) identified above, the `--sync` flag should be used when issuing godcr commands on cli e.g. `godcr send --sync` or `godcr --sync send`. Alternatively, you can set `sync=1` or `sync=true` in `godcr.conf` to always perform a blockchain sync before performing any wallet operation.
- Also, `godcr --sync` can be run alone, without any command to perform a blockchain sync at any time.
- The above concerns do not apply if godcr is _properly_ configured to perform wallet operations using dcrwallet rpc.

### Send and send custom
- May get `wallet.NetworkBackend: Decred network is unreachable` error when using dcrlibwallet; run the command with `--sync` to successfully send funds
- When using dcrlibwallet, successful send transactions do not get published to the blockchain immediately.
The sending account is debited, but the recipient doesn't get the funds until after a while.
There's no solution for this at this time.
It is possible (but not confirmed) that running `godcr --sync` after a send transaction could ensure the transaction gets published to the blockchain.
