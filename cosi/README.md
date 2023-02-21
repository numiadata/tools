# Cosi

Cosi is a cosmos reindexer. It is built to work with many dbs and in the future RPCs.

## How to use cosi

```bash
PROJECT_ID=[project_id] TOPIC=[topic] CHAIN_ID=[chainID] GOOGLE_APPLICATION_CREDENTIALS=[path_to_key] cosi state [state_height] [end_height] [path_to_db] [database_backend (goleveldb or pebbledb)]
```

There are other commands that can be used with cosi as well.

```bash
Usage:
  cosi [command]

Available Commands:
  rpc         reindex via the rpc from a start height to an end height
  kv          reindex via the db from a start height to an end height, note this only works for txs currently
  state       reindex via the state db from a start height to an end height, note this only works for txs currently
  base        Get the base and highest height of the db
  compact     compact the diskspace of the state db
  help        Help about any command

Flags:
  -h, --help     help for cosi
      --unsafe   dont wait on delivery confirmation (default true)

Use "cosi [command] --help" for more information about a command.
```
