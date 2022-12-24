# Pubsub

This package implements google pubsub for usage in an application. 

## How to use this package?

To use this package, first there must be a fork of tendermint that imports it. Once there is a fork that adds this package to the indexing level then a replace tag is required in the applications binary `go.mod`. 

To do this:

```bash
# clone the fork of tendermint at the root level 
$ git clone https://github.com/numiadata/tendermint.git -b v0.34.x

# cd into application directory
cd [application]

# Remove a replace directive.
$ go mod edit -dropreplace google.golang.org/grpc

# Add a replace directive using the fork of tendermint
$ go mod edit -replace github.com/tendermint/tendermint=../tendermint

$ go mod tidy
```

If any errors show up after running go mod tidy some other steps may be needed but it is unknown what those will be without seeing the errors. 


After running this step, we need to set the correct values in the config.toml. 

```toml
#######################################################
###   Transaction Indexer Configuration Options     ###
#######################################################
[tx_index]

# What indexer to use for transactions
#
# The application will set which txs to index. In some cases a node operator will be able
# to decide which txs to index based on configuration set in the application.
#
# Options:
#   1) "null"
#   2) "kv" (default) - the simplest possible indexer, backed by key-value storage (defaults to levelDB; see DBBackend).
# 		- When "kv" is chosen "tx.height" and "tx.hash" will always be indexed.
indexer = "pubsub"
# The Google Cloud Pubsub project ID. Note, operators must ensure the
# GOOGLE_APPLICATION_CREDENTIALS environment variable is set to the location of
# their creds file.
pubsub-project-id = "[insert project id]"

# The Google Cloud Pubsub topic. If the topic does not exist, it will be created.
pubsub-topic = "[insert chainID or topic name here]"
```

After this in order to run a key.json is needed passed as an environment variable.

Reachout to a team member to obtain the key.json.

```bash 
GOOGLE_APPLICATION_CREDENTIALS='/root/key.json' [binary_name] start
```
