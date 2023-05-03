# TX-BOT <!-- omit in toc -->

A tool to execute native Pocket transactions and call Pocket RPCs.

- [Quick Start](#quick-start)
  - [Binary](#binary)
  - [Docker](#docker)
    - [Pulling](#pulling)
    - [Building/Running](#buildingrunning)
- [Usage Options](#usage-options)
  - [Modes of Operation](#modes-of-operation)
  - [Transactions](#transactions)
  - [Supported RPCs](#supported-rpcs)
- [Configuration](#configuration)
- [\[TODO\] Client Generation](#todo-client-generation)
- [Future Work (TODO)](#future-work-todo)

**Too Long; Just Watch (tl;jw)**

If you spin up a [LocalNet](https://github.com/pokt-network/pocket-core/blob/staging/doc/guides/localnet.md), you can either use burst mode or send a single transaction as seen below.

https://user-images.githubusercontent.com/1892194/235985497-82969f83-b2d1-4a3e-ab4d-26a033153713.mov


## Quick Start

### Binary

`make build`

### Docker

#### Pulling

`$ docker-compose up`

#### Building/Running

`$ docker-compose build && docker-compose up`

## Usage Options

### Modes of Operation

- **timerMode**: Execute a random transaction or request every N seconds.
- **randomMode**: Execute a random transaction or request on every key press.
- **selectMode**: Print a list of all the transactions or requests available and trigger one with a key press.
- **manualMode**: Pre-configure a specific transaction or request that can be triggered with a key press.
- **burstMode**: Pre-configure a specific transaction or request that will send N requests separated by M seconds on every key press.

### Transactions

- App Transactions
  - **appStake**: Stake an application for a random set of chains with a random amount for one of the pre-configured keys which was randomly selected.
  - **appUnstake**: Unstake an application for one of the pre-configured keys which was randomly selected.
- Node Transactions
  - **nodeSend**: Send a random amount from one randomly selected address to another.
  - **nodeStake**: Stake a validator for a random set of chains with a random amount for one of the pre-configured keys which was randomly selected.
  - **nodeUnstake**: Unstake a validator node for one of the pre-configured keys which was randomly selected.
  - **nodeUnjail**: Send a node unjail transaction for a randomly selected address.

### Supported RPCs

- **QueryHeight**: Query the height of the Pocket blockchain.
- **relayEthHeight**: Query the height of the Ethereum blockchain by calling `eth_blockNumber`.
- **relayPolyHeight**: Query the height of the Polygon blockchain by calling `eth_blockNumber`.

## Configuration

See the [config](config.json) for an example configuration that works with the `pokt-net-dev-tm`stack in [pocket-e2e-stack](https://github.com/pokt-network/pocket-e2e-stack).

## [TODO] Client Generation

## Future Work (TODO)

[ ] Migrate `tx.go` to use the RPC Client from `req.go`.
[ ] Rename repo to something else since it makes both RPC requests and transactions.
[ ] Pair the configurations for `pocket_endpoint` and `servicer_private_key` so we can distribute relay load.
[ ] Add support for other Ethereum [JSON-RPC requests](https://infura.io/docs/ethereum/json-rpc) APIs.
[ ] Consider consolidating `node_private_keys` and `app_private_keys` in `config.json` and keeping context on which Nodes have made an app stake.
[ ] Add support for other Harmony [JSON-RPC](https://docs.harmony.one/home/developers/api/methods) APIs.
[ ] Consider adding request analytics similar to [PRLTS](https://github.com/pokt-network/prlts).
[ ] Consider using an alternative [Go client generation](https://gist.github.com/craigmurray1120/8e87d88a076d49ec9c43636a313cfa66) or use the Go SDK to be build by PNF in 2022.
[ ] Expand relay chains and functions for each chain.
[ ] Consider having a simple react app as a UI to trigger these transactions or requests.
[ ] Make it easier to send requests to non localnet environments
