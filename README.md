## etherscription - http service supporting notification and subscription

## Idea

### Basic concept

The idea behind this project is to provide a simple http service that allows to subscribe to some ethereum transactions related to a particular address.
All transactions related to subscribed addresses are stored in ephemeral memory storage for time that is no longer than configured retention period.

### The data source

The data source for this project is ethereum node exposing JSON-RPS API. Such a node can be passed as a runtime argument for application.

### The API

This service exposes three endpoints:
- `POST /api/subscribe` - allows subscribing by providing an address `{"address": "0x1234"}`. The address should be a valid ethereum address.
- `GET /api/new-transactions/<address>` - returns all transactions related to subscribed addresses. It is worth mentioning that fetched transactions are removed from storage.
- `GET /api/current-block` - returns id of the latest parse ethereum block.

### How it works under the hood

The service utilizes json-rpc API. Just after starting, it creates a filter for new blocks that appeared in the network using `eth_newFilter`.
Having such a filter, it polls new log entries by `eth_getFilterChanges` method. When a new block appears, it is parsed to extract transaction hashes. 
Then, each transaction is fetched using `eth_getTransactionByHash` method. Fetched transactions are later on mapped to a simplified version and `Value` field is translated from hash encoding to more human-readable format.
Such transactions are stored (if someone subscribed for them OR if service is run with a special flag that allows to store all incoming transactions) in memory storage and are available for fetching

### caveats

You have bear in mind that `eth_newFilter` method is not provided by many public services that overlay ehtereum nodes.
For example, below request
```bash
NODE=https://ethereum-rpc.publicnode.com
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newFilter","params":[{}],"id":2222}' $NODE -s | jq
```

will lead to the following response:

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32701,
    "message": "Please specify an address in your request or, to remove restrictions, order a dedicated full node here: https://www.allnodes.com/eth/host"
  },
  "id": 2222
}
```
To overcome this issue, you can run your own node or use a service that provides such a method.

## Working with service

### Running the service

```
make NODE=YOUR-NODE-ADDRESS-HERE run
```
After running the command, the service will be available under `http://localhost:8888` address and some logs should be printed on stdout:

```
6:42PM WRN Configuration file not found at configuration.yaml, using default path ./configuration/config.yaml module=config
6:42PM INF loaded config
storage:
  retention: 300s
  clean_interval: 3s
  store_all_transactions: true
rpc:
  timeout: 5s
  interval: 3s
  too_many_requests_delay: 500ms
6:42PM INF Filter 0xb80300000000000007848e5cbb31e8fb created module=etherum
6:42PM INF Cleaned 0 entries, left 0 duration=0.017675 module=memory
6:42PM INF Fetched 0 uniq transactions module=etherum
6:42PM INF Cleaned 0 entries, left 0 duration=0.003852 module=memory
6:42PM INF Fetched 168 uniq transactions module=etherum
```

If you with to change the listening port just pass proper PORT variable to make command:

```
make NODE=YOUR-NODE-ADDRESS-HERE PORT=8080 run
```

### Building binaries

```
make build
```

After building, you can run the service by executing:

```
./bin/eherception --node YOUR-NODE-ADDRESS-HERE --port 8080
```

### configuration

The default configuration is stored in `./configuration/config.yaml`. 
You can overwrite it by creating `configuratio.yaml` file in the same place from what you run the service.

Below you can find description of available configuration options:

```yaml
storage:
  retention: 300s # how long transactions are stored in memory
  clean_interval: 3s # how often the goroutine responsible for cleaning is run
  store_all_transactions: true # whether to store all incoming transactions or only those related to subscribed addresses
rpc:
  timeout: 5s # timeout for rpc requests
  interval: 3s # how often the service polls for new transactions
  too_many_requests_delay: 500ms # delay between requests when too many requests are sent to the node
``` 

### Configuring the service

## Development

The project is written in Go. It uses go modules for dependency management. To run tests, you can use `make test` command.

## Further development

There are a lot of things to do to make this project more production ready. Some of them are:
- [ ] Add unit tests to cover more business logic
- [ ] Add integration tests to verify the whole system
- [ ] Add retry logic to handle connection problems with ethereum node
- [ ] Add an ability to easily change logging level & format (code is already prepared for that)
- [ ] Add metrics and expose them via prometheus
- [ ] Switch storage to something durable to avoid loosing data on restart
- [ ] Add parallelism where it makes sense