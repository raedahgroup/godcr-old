# Godcr Status
The current state of godcr is reported below.
Working features, known bugs and issues, work-in-progress features are all listed.

## Features
The list of features available on godcr is maintained below.
Some features are either only partially implemented (:ballot_box_with_check:) or not implemented at all (:x:) on some [interfaces](#interfaces)

| Feature | :white_check_mark: | :ballot_box_with_check: | :x: | Next Steps |
|---|---|---|---|---|
| create wallet | cli, http | terminal | nuklear | Allow creating multiple wallets, even if wallet already exists<br><br>Ask for network type when creating wallet<br><br>Seed display confirmation should follow same pattern as dcrandroid |
| sync blockchain (spv) | cli, http | terminal | nuklear |
| sync blockchain (rpc) | | | all | Add this feature to all interfaces |
| balance | cli, http | nuklear | terminal |
| receive | cli, http, nuklear | | terminal |
| send funds (simple) | cli, http | | nuklear, terminal |
| send funds (custom) | cli | http | nuklear, terminal |
| history | cli, http | | nuklear, terminal |
| tx detail | cli, http | | nuklear, terminal |
| stake info | cli, http | | nuklear, terminal |
| purchase ticket(s) | cli, http | | nuklear, terminal |

## Interfaces
Godcr can run in any of the following interface modes:

#### Cli (`godcr` or `godcr --mode=cli`)
![cli interface screenshot](https://user-images.githubusercontent.com/18400051/52160314-973efd80-26b3-11e9-9eed-7ba0b08f04f4.png)

#### Terminal GUI app (`godcr --mode=terminal`)
![terminal mode screenshot](https://user-images.githubusercontent.com/18400051/52159638-5fca5400-26a7-11e9-877b-54c5f092fbe1.png)

#### Http web app (`godcr --mode=http`)
![http interface screenshot](https://user-images.githubusercontent.com/18400051/52159613-019d7100-26a7-11e9-9cfc-8d044d3468f7.png)

#### Native desktop app (`godcr --mode=nuklear`)
![nuklear mode screenshot](https://user-images.githubusercontent.com/18400051/52159667-d49d8e00-26a7-11e9-9f5f-ba99cb33b4db.png)

## Known Issues
#### Sync blockchain (cli only)
- Unlike in other interfaces, the cli interface does not automatically sync the blockchain before performing wallet operations.
This may lead to inaccuracy of displayed information or complete inability to perform certain wallet operations such as `send` and `purchasetickets`
- To circumvent the issue(s) identified above, the `--sync` flag should be used when issuing godcr commands on cli e.g. `godcr send --sync` or `godcr --sync send`. Alternatively, you can set `sync=1` or `sync=true` in `godcr.conf` to always perform a blockchain sync before performing any wallet operation.
- Also, `godcr --sync` can be run alone, without any command to perform a blockchain sync at any time.
- The above concerns do not apply if godcr is _properly_ configured to perform wallet operations using dcrwallet rpc.

#### Send and send custom (cli only)
- May get `wallet.NetworkBackend: Decred network is unreachable` error when using dcrlibwallet; run the command with `--sync` to successfully send funds
- When using dcrlibwallet, successful send transactions do not get published to the blockchain immediately.
The sending account is debited, but the recipient doesn't get the funds until after a while.
There's no solution for this at this time.
It is possible (but not confirmed) that running `godcr --sync` after a send transaction could ensure the transaction gets published to the blockchain.
